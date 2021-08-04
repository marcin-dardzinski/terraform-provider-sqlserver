package sql

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"time"
)

// adapted from
// https://github.com/Azure/azure-sdk-for-go/blob/fc3a98b94fb25364c8e2799a41adfc369b306166/sdk/azidentity/azure_cli_credential.go#L91
// to include tenant awareness

const timeoutCLIRequest = 10000 * time.Millisecond

type CredentialUnavailableError struct {
	// CredentialType holds the name of the credential that is unavailable
	credentialType string
	// Message contains the reason why the credential is unavailable
	message string
}

func (e *CredentialUnavailableError) Error() string {
	return e.credentialType + ": " + e.message
}

func tenantAwareAzureCLITokenProvider(tenantId, subscriptionId string) func(ctx context.Context, resource string) ([]byte, error) {
	return func(ctx context.Context, resource string) ([]byte, error) {
		// This is the path that a developer can set to tell this class what the install path for Azure CLI is.
		const azureCLIPath = "AZURE_CLI_PATH"

		// The default install paths are used to find Azure CLI. This is for security, so that any path in the calling program's Path environment is not used to execute Azure CLI.
		azureCLIDefaultPathWindows := fmt.Sprintf("%s\\Microsoft SDKs\\Azure\\CLI2\\wbin; %s\\Microsoft SDKs\\Azure\\CLI2\\wbin", os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"))

		// Default path for non-Windows.
		const azureCLIDefaultPath = "/bin:/sbin:/usr/bin:/usr/local/bin"

		// Validate resource, since it gets sent as a command line argument to Azure CLI
		const invalidResourceErrorTemplate = "resource %s is not in expected format. Only alphanumeric characters, [dot], [colon], [hyphen], and [forward slash] are allowed"
		match, err := regexp.MatchString("^[0-9a-zA-Z-.:/]+$", resource)
		if err != nil {
			return nil, err
		}
		if !match {
			return nil, fmt.Errorf(invalidResourceErrorTemplate, resource)
		}

		ctx, cancel := context.WithTimeout(ctx, timeoutCLIRequest)
		defer cancel()

		// Execute Azure CLI to get token
		var cliCmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cliCmd = exec.CommandContext(ctx, fmt.Sprintf("%s\\system32\\cmd.exe", os.Getenv("windir")))
			cliCmd.Env = os.Environ()
			cliCmd.Env = append(cliCmd.Env, fmt.Sprintf("PATH=%s;%s", os.Getenv(azureCLIPath), azureCLIDefaultPathWindows))
			cliCmd.Args = append(cliCmd.Args, "/c", "az")
		} else {
			cliCmd = exec.CommandContext(ctx, "az")
			cliCmd.Env = os.Environ()
			cliCmd.Env = append(cliCmd.Env, fmt.Sprintf("PATH=%s:%s", os.Getenv(azureCLIPath), azureCLIDefaultPath))
		}
		cliCmd.Args = append(cliCmd.Args, "account", "get-access-token", "-o", "json", "--resource", resource)

		if tenantId != "" {
			cliCmd.Args = append(cliCmd.Args, "--tenant", tenantId)
		}

		if subscriptionId != "" {
			cliCmd.Args = append(cliCmd.Args, "--subscription", subscriptionId)
		}

		var stderr bytes.Buffer
		cliCmd.Stderr = &stderr

		output, err := cliCmd.Output()
		if err != nil {
			msg := stderr.String()
			if msg == "" {
				// if there's no output in stderr report the error message instead
				msg = err.Error()
			}
			return nil, &CredentialUnavailableError{credentialType: "Azure CLI Credential", message: msg}
		}

		return output, nil
	}
}
