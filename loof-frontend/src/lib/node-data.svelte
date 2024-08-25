<script lang="ts" context="module">
    import type { RecordModel } from "pocketbase";
	import PocketBase, { ClientResponseError } from "pocketbase"

	export enum NodeState {
		moving,
		editing,
		complete
	}

	export class NodeData {
		id: string | null;

		state: NodeState;

		authorId: string;
		previousNodeId: string;

		x: number | null;
		y: number | null;

		text: string;

		constructor(authorId: string, previousNodeId: string, x: number | null, y: number | null) {
			this.id = null;

			this.state = NodeState.moving;

			this.authorId = authorId;
			this.previousNodeId = previousNodeId;

			this.x = x;
			this.y = y;

			this.text = "";
		}
	}

	export function fromDatabase(databaseRecord: RecordModel): NodeData {
			return {
				id: databaseRecord.id,

				state: NodeState.complete,

				authorId: databaseRecord.author,
				previousNodeId: databaseRecord.previous_node,

				x: databaseRecord.x_position,
				y: databaseRecord.y_position,

				text: databaseRecord.text
			}
		}
	
	export function toDatabase(node: NodeData): any {
		return {
				author: node.authorId,
				previous_node: node.previousNodeId,

				x_position: node.x,
				y_position: node.y,

				text: node.text
			}
	}

	export async function checkIfUserIsAuthor(authorId: string | null, pocketbase: PocketBase): Promise<boolean> {
		if (authorId == null) {
			return false;
		} else {
			let authorRecord = await pocketbase.collection('authors').getFirstListItem('author_id?="' + authorId + '"&&' + 'origin="user"')
									.catch(err => {if ((err instanceof ClientResponseError) && (err.status) == 404) { return null }});
				

			if (authorRecord?.id != null) {
				return true;
			} else {
				return false;
			}
		}
		
	}

	export async function getDisplayUsername(node: NodeData, pocketbase: PocketBase): Promise<string | null> {
		let authorRecord = await pocketbase.collection('authors').getOne(node.authorId);

		if (authorRecord.id != null) {
			if (authorRecord.type == "human") {
				let userRecord = await pocketbase.collection('users').getOne(authorRecord.authorId); 
				return userRecord.username;
			} else {
				let userRecord = await pocketbase.collection('robots').getOne(authorRecord.authorId); 
				return userRecord.name;
			}
		} else {
			return null;
		}
	}

</script>
