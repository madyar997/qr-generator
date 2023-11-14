package usecase

import (
	"context"
	"fmt"
	ssoClient "github.com/madyar997/user-client/client/grpc"
	"github.com/opentracing/opentracing-go"
	spanLog "github.com/opentracing/opentracing-go/log"
	qrcode "github.com/skip2/go-qrcode"
	"log"
	"net/http"
	"strconv"
)

// QrUseCase -.
type QrUseCase struct {
	httpCli *http.Client
	ssoCli  *ssoClient.Client
}

// New -.
func NewQrUseCase(httpCli *http.Client, client *ssoClient.Client) *QrUseCase {
	return &QrUseCase{
		httpCli: httpCli,
		ssoCli:  client,
	}
}

func (uq *QrUseCase) Me(ctx context.Context, userID int) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "qr generator use case")
	defer span.Finish()

	span.LogFields(
		spanLog.String("user_id", strconv.Itoa(userID)),
	)

	user, err := uq.ssoCli.GetUserByID(ctx, int32(userID))
	if err != nil {
		log.Printf("could not get user from sso %s", err)
		return nil, err
	}

	var png []byte
	png, err = qrcode.Encode(fmt.Sprintf("https://example.com/%s", user.Email), qrcode.Medium, 256)
	if err != nil {
		log.Printf("could not create qr code %s", err)
		return nil, err
	}
	return png, nil
}
