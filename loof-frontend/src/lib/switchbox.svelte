<script lang="ts">
	let { x_position=0, y_position=0, on_write_selection, on_generate_selection }: { x_position: number, y_position: number, on_write_selection: () => void, on_generate_selection: () => void } = $props();
</script>

<div class="switchbox-container" style="top: {y_position.toString()}px; left: {x_position.toString()}px;">
	<ul>
		<li><button onclick={on_write_selection}>Write (as a human)</button></li>
		<li><button onclick={on_generate_selection}>Generate (using llama-3.1-405b-base)</button></li>
	</ul>
</div>

<style>
	.switchbox-container {
		position: absolute; 

    	width: 300px;

		border: 5px solid rgb(66, 66, 66);
		box-sizing: border-box;

		background-color: rgb(236, 232, 223);
	}

	ul {
		list-style: none;

		padding: 0;
	}

	button {
		background-color: inherit;

		border: 0;

  		padding: 0;
		padding-left: 5px;

		cursor: pointer;
	}
</style>