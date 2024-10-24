<div class="modal-container">
	<slot />
</div>

<style>
	.modal-container {
		overflow: hidden;
		position: fixed;

		inset: 0;

		width: fit-content;
		height: fit-content;

		margin: auto;
		padding: 30px;

		border: 10px solid black;
		background-color: white;

		max-width: 30%;
		max-height: 50%;
	}

	@media (max-width: 768px) {
        .modal-container { max-width: 100%; max-height: 100%; height: 100%; width: 100%; box-sizing: border-box; }
    }
</style>