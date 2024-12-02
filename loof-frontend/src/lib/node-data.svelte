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
		isLocal: boolean;
		fromRobot: boolean;

		authorId: string;
		previousNodeId: string;

		x: number | null;
		y: number | null;

		text: string;

		constructor(authorId: string, previousNodeId: string, x: number | null, y: number | null) {
			this.id = null;

			this.state = NodeState.moving;
			this.isLocal = true;
			this.fromRobot = false;

			this.authorId = authorId;
			this.previousNodeId = previousNodeId;

			this.x = x;
			this.y = y;

			this.text = "";
		}
	}

	function convertDatabaseStateToInternalState(dbState: string): NodeState | null {
		switch (dbState) {
			case "moving":
				return NodeState.moving;
			case "pending":
				return NodeState.editing;
			case "editing":
				return NodeState.editing;
			case "completed":
				return NodeState.complete;
			default:
				return null;
		}
	}

	function convertInternalStateToDatabaseState(state: NodeState, afterPending: boolean): string {
		switch (state) {
			case NodeState.moving:
				return "moving";
			case NodeState.editing:
				if (afterPending == true) {
					return "editing";
				} else {
					return "pending";
				}
			case NodeState.complete:
				return "completed";
		}
	}

	export function fromDatabase(databaseRecord: RecordModel, isLocal: boolean): NodeData {
			return {
				id: databaseRecord.id,

				state: convertDatabaseStateToInternalState(databaseRecord.state)!,
				isLocal: isLocal,
				fromRobot: (databaseRecord.expand!.author.origin == "robot") ? true : false,

				authorId: databaseRecord.author,
				previousNodeId: databaseRecord.previous_node,

				x: databaseRecord.x_position,
				y: databaseRecord.y_position,

				text: databaseRecord.text
			}
		}
	
	export function toDatabase(node: NodeData, ignoreState: boolean): any {
		let incompleteNode: any = {
			author: node.authorId,
			previous_node: node.previousNodeId,

			x_position: node.x,
			y_position: node.y,

			text: node.text
		}

		if (ignoreState == false) {
			incompleteNode.state = convertInternalStateToDatabaseState(node.state, false)
		}

		return incompleteNode;
	}

	export async function checkIfUserIsAuthor(authorId: string | null, pocketbase: PocketBase): Promise<boolean> {
		if (authorId == null) {
			return false;
		} else {
			let authorRecord = await pocketbase.collection('authors').getFirstListItem('author_id?="' + authorId + '"&&' + 'origin="human"')
									.catch(err => {if ((err instanceof ClientResponseError) && (err.status) == 404) { return null }});
				

			if (authorRecord?.id != null) {
				return true;
			} else {
				return false;
			}
		}	
	}

	export async function getUserAuthor(authorId: string | null, pocketbase: PocketBase): Promise<string | null> {
		if (authorId == null) {
			return null;
		} else {
			let authorRecord = await pocketbase.collection('authors').getFirstListItem('author_id?="' + authorId + '"&&' + 'origin="human"', { requestKey: null })
									.catch(err => {if ((err instanceof ClientResponseError) && (err.status) == 404) { return null }});

			return authorRecord?.id ?? null;
		}	
	}

	export async function checkIfUserIsAuthorOfNode(authorId: string | null, pocketbase: PocketBase): Promise<boolean> {
		if (authorId == null) {
			return false;
		} else {
			let authorRecord = await pocketbase.collection('authors').getOne(authorId, { requestKey: null })
									.catch(err => {if ((err instanceof ClientResponseError) && (err.status) == 404) { return null }});
			

			if (authorRecord?.id == authorId) {
				return true;
			} else {
				return false;
			}
		}	
	}

	// TODO: Make this function more efficient
	export async function getDisplayUsername(node: NodeData, pocketbase: PocketBase): Promise<string | null> {
		if (node.authorId != null) {
			if (node.fromRobot == false) {
				let authorRecord = await pocketbase.collection('authors').getOne(node.authorId, { requestKey: null });

				let userRecord = await pocketbase.collection('users').getOne(authorRecord.author_id, { requestKey: null }); 
				return userRecord.username;
			} else {
				let authorRecord = await pocketbase.collection('authors').getOne(node.authorId, { requestKey: null });

				let robotRecord = await pocketbase.collection('robots').getOne(authorRecord.author_id, { requestKey: null }); 
				return robotRecord.model;
			}
		} else {
			return null;
		}
	}

	export async function getQuotationUsername(authorId: string, pocketbase: PocketBase) : Promise<string> {
		let authorRecord = await pocketbase.collection('authors').getOne(authorId, { requestKey: null });

		if (authorRecord.origin == "human") {
			let userRecord = await pocketbase.collection('users').getOne(authorRecord.author_id, { requestKey: null }); 
			return userRecord.username;
		} else {
			let robotRecord = await pocketbase.collection('robots').getOne(authorRecord.author_id, { requestKey: null }); 
			return robotRecord.model;
		}
	}

</script>
