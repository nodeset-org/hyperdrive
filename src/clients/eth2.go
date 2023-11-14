package eth2_client

import (
    "fmt"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
)

type JSONRPCRequest struct {
    JSONRPC string      `json:"jsonrpc"`
    Method  string      `json:"method"`
    Params  interface{} `json:"params"`
    ID      int         `json:"id"`	
}

func CheckStatus(clientURL string) (bool, error) {
    fmt.Println("Connecting to ETH2 Client...")

    // Get current blockNumber as a status check
    requestBody, err := json.Marshal(JSONRPCRequest{
        JSONRPC: "2.0",
        Method:  "eth_blockNumber",
        Params:  []interface{}{},
        ID:      1,
    })

    if err != nil {
        return false, err
    }

    req, err := http.NewRequest("POST", clientURL, bytes.NewBuffer(requestBody))

    if err != nil {
        return false, err
    }

    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)

    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return false, err
    }
    fmt.Printf("Response body: %s\n", body)
    return true, nil
}
