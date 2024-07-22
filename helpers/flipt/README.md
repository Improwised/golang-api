# Golang Boilerplate with Flipt Integration

This project integrates [Flipt](https://flipt.io) for feature flag management.

## What is Flipt?

Flipt is an open-source feature flag management tool that allows you to enable or disable features in your application without deploying new code. It provides a simple way to control the availability of features to your users, perform A/B testing, and gradually roll out new features.

## Why Use Flipt?

- **Feature Management**: Easily turn features on or off based on various conditions.
- **A/B Testing**: Run experiments to determine which features perform better.
- **Gradual Rollout**: Gradually release features to a subset of users to minimize risk.
- **Decoupled Releases**: Deploy code without enabling the feature, reducing the risk of breaking changes.

## Configuration

1. Set the `FLIPT_ENABLED` environment variable to `true` to enable Flipt support.
2. Set the `FLIPT_HOST` (e.g., `localhost`) and `FLIPT_PORT` (e.g., `9000`) environment variables to communicate with the Flipt backend.
3. For setting up the ```Flipt``` service with other services, you have to execute ```docker-compose --profile flipt up```.
4. Now you can access flipt UI on [http://localhost:8081](http://localhost:8081)

## Changes Made to Integrate Flipt

1. **Initialize Flipt Client on Application Startup:**

    In `api.go`, add the following code to initialize the Flipt client:

    ```go
    // Initialize Flipt client for flipt functionality
    err = helpers.InitFliptClient()
    if err != nil {
        logger.Error("Error while initializing client", zap.Error(err))
        if err.Error() != "flipt is not enabled" {
            return err
        }
    }
    ```

2. **Helper Functions to Communicate with Flipt Backend:**

    The `helpers` package includes two methods for interacting with Flipt:

    - **GetBooleanFlag**: Retrieves a boolean flag based on the flag key.

        ```go
        type BooleanFlagResponse struct {
            Key         string `json:"key"`
            Name        string `json:"name"`
            Description string `json:"description"`
            Enabled     bool   `json:"enabled,omitempty"`
        }

        func GetBooleanFlag(flagKey string) (BooleanFlagResponse, error) {
            // Implementation here
        }
        ```

    - **GetVariantFlag**: Retrieves a variant flag based on the flag key, entity ID, and context map.

        ```go
        type Context struct {
            Key   string `json:"key"`
            Value string `json:"value"`
        }

        type VariantFlagResponse struct {
            RequestContext Context `json:"request_context"`
            Match          bool    `json:"match,omitempty"`
            FlagKey        string  `json:"flag_key"`
            SegmentKey     string  `json:"segment_key,omitempty"`
            Value          string  `json:"value,omitempty"`
        }

        func GetVariantFlag(flagKey string, entityId string, contextMap map[string]string) (VariantFlagResponse, error) {
            // Implementation here
        }
        ```

## How to Use Flipt in Your Code

### Use Boolean Flag

1. Open [Flipt Flags](http://localhost:8081/#/flags) and create a new flag:
    - Enter the name and key for the flag (e.g., `advertisement`).
    - Set the type to boolean, add a description, enable the flag, and click create.

2. Use the created flag in your code:

    ```go
    advertisementFlag, err := helpers.GetBooleanFlag("advertisement")
    if err != nil {
        hc.logger.Error("error while checking flag from Flipt", zap.Error(err))
        return err
    }

    if advertisementFlag.Enabled {
        // Your logic here
    }
    ```

### Use Variant Flag

1. Open [Flipt Flags](http://localhost:8081/#/flags) and create a new flag:
    - Enter the name and key for the flag (e.g., `color`).
    - Add a description, enable the flag, and click create.

2. Create variants for the flag:
    - Add variants (e.g., `orange`, `green`).

3. Create a segment:
    - Add constraints based on your requirements (e.g., `country` equals `ind`).

4. Create a rule:
    - Assign the segment to the flag and specify the variant to return when the segment matches.

5. Use the created flag in your code:

    ```go
    variantFlag, err := helpers.GetVarientFlag("color", "1234", map[string]string{"country": "ind"})
    if err != nil {
        hc.logger.Error("error while checking flag from Flipt", zap.Error(err))
        return err
    }

    if variantFlag.Value == "orange" {
        // Your logic here
    } else if variantFlag.Value == "green" {
        // Your logic here
    }
    ```

    - `GetVariantFlag` method:
        ```go
        flagKey = "color"
        entityId = "1234"
        contextMap = map[string]string{"country": "ind"}
        ```

    Here, `country: ind` is the constraint and `entityId: 1234` is the entity ID.
