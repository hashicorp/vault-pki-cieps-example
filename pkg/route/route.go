package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/vault-pki-cieps-example/internal/business"

	"github.com/hashicorp/vault/sdk/helper/certutil"
)

type (
	httpHandlerFunc    func(resp http.ResponseWriter, req *http.Request)
	ErrableHandlerFunc func(resp http.ResponseWriter, req *http.Request) error
)

func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/evaluate", ErrLogger(EvaluatePolicy))
}

var FatalError = errors.New("internal fatal error")

func ErrLogger(wrapped ErrableHandlerFunc) httpHandlerFunc {
	return func(respWriter http.ResponseWriter, req *http.Request) {
		if err := wrapped(respWriter, req); err != nil {
			log.Printf("[%v@%v] failed: %v", req.Method, req.URL.String(), err)
			errMsg := err.Error()
			errStatus := http.StatusInternalServerError
			if !strings.Contains(err.Error(), FatalError.Error()) {
				resp := &certutil.CIEPSResponse{
					Error: err.Error(),
				}
				respBytes, jsonErr := json.Marshal(resp)
				if jsonErr != nil {
					// Fallback to standard text encoding and hope it doesn't
					// error further.
					errMsg = `{"error":"` + err.Error() + `",warnings:["failed json encoding error response: ` + jsonErr.Error() + `"]}`
				} else {
					errMsg = string(respBytes)
				}
				errStatus = http.StatusOK
			}
			http.Error(respWriter, errMsg, errStatus)
		}
	}
}

func EvaluatePolicy(respWriter http.ResponseWriter, req *http.Request) error {
	policyReq, err := ParsePolicyRequest(req)
	if err != nil {
		return fmt.Errorf("error evaluating policy: %w", err)
	}

	policyResp, err := business.Evaluate(policyReq)
	policyResp.UUID = policyReq.UUID

	logResp(policyReq, policyResp)

	return SendPolicyResp(respWriter, policyResp)
}

func ParsePolicyRequest(req *http.Request) (*certutil.CIEPSRequest, error) {
	if req.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("%v: invalid content type (%v); expected application/json", FatalError, req.Header.Get("Content-Type"))
	}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	var request certutil.CIEPSRequest
	if err := decoder.Decode(&request); err != nil {
		return nil, fmt.Errorf("%v: invalid request from Vault: failed to decode JSON request body data: %w", FatalError, err)
	}

	if err := request.ParseUserCSR(); err != nil {
		return nil, fmt.Errorf("%v: invalid request from Vault: expected valid CSR on request: %w", FatalError, err)
	}

	return &request, nil
}

func logResp(policyReq *certutil.CIEPSRequest, policyResp *certutil.CIEPSResponse) {
	reqBytes, err := json.Marshal(policyReq)
	req := string(reqBytes)
	if err != nil {
		req = fmt.Sprintf("failed marshaling request: [err: %v / req: %v]", err, policyReq)
	}

	respBytes, err := json.Marshal(policyResp)
	resp := string(respBytes)
	if err != nil {
		resp = fmt.Sprintf("failed marshaling response: [err: %v / resp: %v]", err, policyResp)
	}

	log.Printf("Evaluated policy request: [req: %v] -> [resp: %v]", req, resp)
}

func SendPolicyResp(respWriter http.ResponseWriter, policyResp *certutil.CIEPSResponse) error {
	respBytes, err := json.Marshal(policyResp)
	if err != nil {
		return fmt.Errorf("%v: failed to marshal response object: %w", FatalError, err)
	}

	// We like the default behavior of writing 200 OK; here, even if the
	// policy evaluation has resulted in an error, we wish to return a
	// Vault-level error and not a panic'ing error.
	for len(respBytes) > 0 {
		written, err := respWriter.Write(respBytes)
		if err != nil {
			return fmt.Errorf("%v: failed writing to response stream: %w", FatalError, err)
		}
		if written == 0 {
			return fmt.Errorf("%v: failed to write any bytes to response stream: 0 of %v", FatalError, len(respBytes))
		}

		respBytes = respBytes[written:]
	}

	return nil
}
