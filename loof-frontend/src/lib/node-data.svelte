<script lang="ts" context="module">
	export enum NodeState {
		moving,
		editing
	}

	export class InConstructionNode {
		state: NodeState;
		x: number | null;
		y: number | null;

		constructor(x: number, y: number) {
			this.state = NodeState.moving;

			this.x = x;
			this.y = y;
		}
	}
</script>
