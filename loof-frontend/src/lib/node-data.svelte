<script lang="ts" context="module">
	export enum NodeState {
		moving,
		editing
	}

	export class InConstructionNode {
		state: NodeState;
		previousNodeId: string;

		x: number | null;
		y: number | null;

		constructor(previousNodeId: string, x: number, y: number) {
			this.state = NodeState.moving;
			this.previousNodeId = previousNodeId;

			this.x = x;
			this.y = y;
		}
	}
</script>
