<script lang="ts">
	let { x_position=0, y_position=0 } = $props();
</script>

<div class="switchbox-container" style="top: {y_position.toString()}px; left: {x_position.toString()}px;">
	<ul>
		<li><button>Write (as a human)</button></li>
		<li><button>Generate (using llama-3.1-405b-base)</button></li>
	</ul>
</div>

<style>
	.switchbox-container {
		position: absolute; 

    	width: 200px;

		border: 5px solid rgb(66, 66, 66);
		box-sizing: border-box;

		background-color: rgb(236, 232, 223);
	}

	button {
		background-color: inherit;

		border: 0;

  		padding: 0;
		padding-left: 5px;
	}
</style>