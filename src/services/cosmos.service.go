package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	types4 "github.com/cosmos/cosmos-sdk/codec/types"
	codec2 "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	types2 "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	signing2 "github.com/cosmos/cosmos-sdk/x/auth/signing"
	tx2 "github.com/cosmos/cosmos-sdk/x/auth/tx"
	types3 "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sirupsen/logrus"
	"github.com/sonntuet1997/medical-chain-utils/cryptography"
	types "github.com/sotanext-team/medical-chain/src/sideservice/src/pb/medichain"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type CosmosServiceI interface {
	GetUser(req *types.QueryGetUserRequest) (*types.QueryGetUserResponse, error)
	GetSharing(req *types.QueryGetSharingRequest) (*types.QueryGetSharingResponse, error)
	GetCheckSharing(req *types.QueryCheckSharingRequest) (*types.QueryCheckSharingResponse, error)
	GetAccount(address string) (*types3.BaseAccount, error)
	PostCreateServiceUser(req *types.MsgCreateServiceUser) (*txtypes.BroadcastTxResponse, error)
}

type CosmosService struct {
	queryClient types.QueryClient
	authClient  types3.QueryClient
	txClient    txtypes.ServiceClient
	ctx         context.Context
	chainId     string
	keyring     keyring.UnsafeKeyring
	cdc         *codec.ProtoCodec
}

func InitCodec() (cdc *codec.ProtoCodec) {
	registry := types4.NewInterfaceRegistry()
	types3.RegisterInterfaces(registry)
	codec2.RegisterInterfaces(registry)
	types.RegisterInterfaces(registry)
	cdc = codec.NewProtoCodec(registry)
	return cdc
}

func NewCosmosService(ctx context.Context, keyring keyring.UnsafeKeyring, chainId string, cosmosEndpoint string) *CosmosService {
	cc, err := grpc.DialContext(ctx, cosmosEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil
	}
	qClient := types.NewQueryClient(cc)
	authClient := types3.NewQueryClient(cc)
	txClient := txtypes.NewServiceClient(cc)

	return &CosmosService{
		queryClient: qClient,
		authClient:  authClient,
		txClient:    txClient,
		ctx:         ctx,
		chainId:     chainId,
		keyring:     keyring,
		cdc:         InitCodec(),
	}
}

func (s *CosmosService) GetUser(req *types.QueryGetUserRequest) (*types.QueryGetUserResponse, error) {
	res, err := s.queryClient.User(s.ctx, req)

	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}
	return res, nil
}

func (s *CosmosService) GetSharing(req *types.QueryGetSharingRequest) (*types.QueryGetSharingResponse, error) {
	res, err := s.queryClient.Sharing(s.ctx, req)

	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}
	return res, nil
}

func (s *CosmosService) GetCheckSharing(req *types.QueryCheckSharingRequest) (*types.QueryCheckSharingResponse, error) {
	res, err := s.queryClient.CheckSharing(s.ctx, req)

	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}
	if res.Sharing.Status != "accepted" {
		return nil, xerrors.New("invalid permission")
	}

	return res, nil
}

func (s *CosmosService) GetAccount(address string) (*types3.BaseAccount, error) {
	req := types3.QueryAccountRequest{Address: address}
	res, err := s.authClient.Account(s.ctx, &req)
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}
	var account = types3.BaseAccount{}

	bz, err := s.cdc.MarshalJSON(res.Account)
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}

	err = json.Unmarshal(bz, &account)
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}

	return &account, nil
}

func (s *CosmosService) PostCreateServiceUser(req *types.MsgCreateServiceUser) (*txtypes.BroadcastTxResponse, error) {
	adminPrivStr, err := s.keyring.UnsafeExportPrivKeyHex("admin")
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}

	adminPrivBz, err := hex.DecodeString(adminPrivStr)
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}
	adminPriv := secp256k1.PrivKey{Key: adminPrivBz}
	privs := []types2.PrivKey{&adminPriv}

	res, err := s.sendTx(req, privs)
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}
	return res, nil
}

func (s *CosmosService) AddAccountFromMnemonic(uid string, mnemonic string) (keyring.Info, error) {
	return s.keyring.NewAccount(uid, mnemonic, "", "", hd.Secp256k1)
}

func (s *CosmosService) ShowAccount(uid string) (keyring.Info, error) {
	return s.keyring.Key(uid)
}

func (s *CosmosService) sendTx(msg sdk.Msg, privs []types2.PrivKey) (*txtypes.BroadcastTxResponse, error) {
	var accs []*types3.BaseAccount

	for _, priv := range privs {
		addr, err := sdk.Bech32ifyAddressBytes("medichain", priv.PubKey().Address().Bytes())

		if err != nil {
			message := status.Convert(err).Message()
			return nil, xerrors.New(message)
		}
		logrus.Info(addr)
		acc, err := s.GetAccount(addr)
		if err != nil {
			message := status.Convert(err).Message()
			return nil, xerrors.New(message)
		}
		logrus.Info(acc)
		accs = append(accs, acc)
	}

	txConfig := tx2.NewTxConfig(s.cdc, tx2.DefaultSignModes)

	txBuilder := txConfig.NewTxBuilder()

	err := txBuilder.SetMsgs(msg)
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}

	txBuilder.SetGasLimit(200000)

	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  txConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: accs[i].Sequence,
		}

		sigsV2 = append(sigsV2, sigV2)
	}

	err = txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}

	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := signing2.SignerData{
			ChainID:       s.chainId,
			AccountNumber: accs[i].AccountNumber,
			Sequence:      accs[i].Sequence,
		}
		sigV2, err := tx.SignWithPrivKey(
			txConfig.SignModeHandler().DefaultMode(), signerData,
			txBuilder, priv, txConfig, accs[i].Sequence)
		if err != nil {
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}

	err = txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}

	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}

	txRes, err := s.txClient.BroadcastTx(
		s.ctx,
		&txtypes.BroadcastTxRequest{
			Mode:    txtypes.BroadcastMode_BROADCAST_MODE_BLOCK,
			TxBytes: txBytes, // Proto-binary of the signed transaction, see previous step.
		},
	)
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}

	return txRes, nil
}

func (s *CosmosService) Sign(msg string, key string) ([]byte, error) {

	//key64 := "ggR2Oy+S8Imls539nCBu+T9Rv8NO9JHTxuFOK75nOpA="
	//logrus.Info(key64)

	keyBz, err := cryptography.ConvertBase64ToBytes(key)
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}

	sig, err := cryptography.SignMessage(msg, keyBz)
	if err != nil {
		message := status.Convert(err).Message()
		return nil, xerrors.New(message)
	}

	logrus.Info(cryptography.ConvertBytesToBase64(sig))

	return sig, nil
}

func (s *CosmosService) Verify(msg string, sig string, pubKey string) (bool, error) {

	//pubKey64 := "A8SmZkaxqfKKjC6JoxQ/0o9WSG0D9clAuxrHniut0l52"
	//logrus.Info(pubKey64)
	pubKeyBz, err := cryptography.ConvertBase64ToBytes(pubKey)
	if err != nil {
		logrus.Error(err)
		message := status.Convert(err).Message()
		return false, xerrors.New(message)
	}

	//sig, err := cryptography.SignMessage("sdfsjflsdj", pubKey64)
	//if err != nil {
	//	message := status.Convert(err).Message()
	//	return nil, xerrors.New(message)
	//}

	sigBz, err := cryptography.ConvertBase64ToBytes(sig)
	if err != nil {
		logrus.Error(err)
		message := status.Convert(err).Message()
		return false, xerrors.New(message)
	}

	verifySig, err := cryptography.VerifySig(msg, sigBz, pubKeyBz)
	if err != nil {
		logrus.Error(err)
		message := status.Convert(err).Message()
		return false, xerrors.New(message)
	}

	logrus.Info(verifySig)

	return verifySig, nil
}
