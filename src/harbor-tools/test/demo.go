package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	authConfig := types.AuthConfig{
		Username: "admin",
		Password: "Harbor12345",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		log.Printf("encode json error:%s", err)
	}
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}


	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := cli.ImagePull(ctx, "harbor.renmaitech.cn/bestriver/best-river-admin:v1.2.4", types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}

	defer out.Close()
	io.Copy(os.Stdout, out)
	fmt.Println(">>>>>>")
}