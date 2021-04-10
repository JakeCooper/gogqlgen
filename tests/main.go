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

	func getAllContainers(ctx context.Context, req *entity.ALlContainersRequest)
*/
