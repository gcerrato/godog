// todo: move to backend package
package llm

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func SendToLLM(ctx context.Context, input string) (string, error) {
	llm, err := openai.New()
    if err != nil {
        log.Fatal(err)
    }
	fmt.Println(input)
    completion, err := llm.Call(ctx, input,
        llms.WithTemperature(0.8),
    )
    if err != nil {
        log.Fatal(err)
    }

	return completion, nil
}
