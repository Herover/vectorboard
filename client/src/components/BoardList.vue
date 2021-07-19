<template>
  <ul class="board-list">
    <li v-for="board in boards" v-bind:key="board._id">
      {{ board.name }}
    </li>
  </ul>
</template>

<script>
export default {
  name: 'BoardList',
  props: {
    boards: Array
  }
}
</script>

<style scoped>

</style>