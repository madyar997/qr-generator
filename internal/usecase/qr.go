package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/madyar997/qr-generator/internal/entity"
	"github.com/madyar997/qr-generator/pkg/jaeger"
	"github.com/opentracing/opentracing-go"
	spanLog "github.com/opentracing/opentracing-go/log"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// QrUseCase -.
type QrUseCase struct {
	httpCli *http.Client
}

// New -.
func NewQrUseCase(httpCli *http.Client) *QrUseCase {
	return &QrUseCase{
		httpCli: httpCli,
	}
}

func (uq *QrUseCase) Me(ctx context.Context, userID int) (*entity.UserInfo, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "qr generator use case")
	defer span.Finish()

	span.LogFields(
		spanLog.String("user_id", strconv.Itoa(userID)),
	)

	url := fmt.Sprintf("http://localhost:8080/api/v1/admin/user/%d", userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error %s", err)
		return nil, err
	}

	err = jaeger.Inject(span, req)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", `application/json`)
	resp, err := uq.httpCli.Do(req)
	if err != nil {
		log.Printf("error %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	var user *entity.UserInfo
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error %s", err)
		return nil, err
	}

	if err = json.Unmarshal(data, &user); err != nil {
		log.Printf("error %s", err)
		return nil, err
	}

	return user, nil
}
