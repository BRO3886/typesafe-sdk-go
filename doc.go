// Package typesafe is the Go SDK for the TypeSafe AI API.
//
// Ask named questions about text or structured state and read back typed
// answers. Set TYPESAFE_API_KEY in the environment, then:
//
//	client, err := typesafe.New()
//	if err != nil {
//		return err
//	}
//	resp, err := client.SystemOne(ctx,
//		"I was charged twice. Please fix this ASAP.",
//		typesafe.Questions{
//			"category": typesafe.Choice{
//				Instructions: "What is this ticket about?",
//				Criteria: map[string]typesafe.Content{
//					"billing": nil, "technical": nil, "other": nil,
//				},
//			},
//		},
//	)
//	if err != nil {
//		return err
//	}
//	fmt.Println(resp.Choices()["category"].Choice)
//
// Every call takes a [context.Context]: cancelling it aborts the request and
// any pending retry, and its deadline bounds the call as a whole, including
// backoff. The per-attempt timeout is configured separately with [WithTimeout].
//
// Learn what TypeSafe can do in the TypeSafe docs at https://docs.typesafe.ai/.
package typesafe
