package main

/*
	In:
	query getMongoCollectionNames(
		$projectId: String!
		$environmentId: String!
		$pluginId: String!
	) {
		mongoCollectionNames(
			projectId: $projectId
			environmentId: $environmentId
			pluginId: $pluginId
		)
	}

	Out:
	func MongoCollectionNames(ctx context.Context, req *entity.MongoCollectionNamesRequest) []string {

	}

	type MongoCollectionNamesRequest struct {
		ProjectID: string
		EnvironmentID: string
		PluginID: string
	}

	In:
	query getAllContainers($limit: Int, $offset: Int) {
		allContainers(limit: $limit, offset: $offset) {
			...ContainerFields
			plugin {
				project {
					id
					name
				}
			}
		}
	}

	Out:

	func getAllContainers(ctx context.Context, req *entity.ALlContainersRequest) *entity.AllContainersResponse {

	}

	type AllContainersRequest struct {
		Limit: Int
		Offset: Int
		GQL: entity.AllContainersGQL {
			ID: bool
			CreatedAt: bool
			Envs: bool
			Plugin: entity.PluginGQL {
				ID: bool
				Name: bool
				Project: entity.ProjectGQL {
					ID: bool
					Name: bool
				}
			}
		}
	}

	type AllContainersResponse struct {
		ID:
		Plugin: Plugin {
			Project: Project {
				ID: string
				Name: string
			}
		}
	}

	In:
	ENUM DeployStatus

	Out:
	const (
		STATUS_BUILDING DeployStatus
	)
*/
