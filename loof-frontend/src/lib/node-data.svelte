<script lang="ts" context="module">
    import type { RecordModel } from "pocketbase";

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

</script>
