<template>
  <div class="about">
    <h1>{{ board ? board.name : 'Loading...' }}</h1>
    <button v-on:click="deleteBoard()">Delete</button><!-- TODO: Only display if we can do this -->
    <input type="checkbox" id="hide-board" v-model="hidden"/><label for="hide-board">Hidden</label>
    <br/>
    <br/>
    <button v-on:click="addText()">Add text</button>
    <br/>
    <svg width="800" height="600" v-if="board">
      <g v-for="item in board.content" :key="item.id">
        <text v-if="item.type == 'text'" :x="item.x" :y="item.y">{{ item.string }}</text>
      </g>
    </svg>
  </div>
</template>
<script>
import { API_BASE, WS_BASE } from '@/config.js'

export default {
  name: 'Board',
  components: {
  },
  data: () => ({
    socket: null,
    board: null,
    hidden: false,
  }),
  watch: {
    hidden: async function(oldVal, newVal) {
      console.log(API_BASE)
      const resp = await fetch(`${API_BASE}/boards/${this.$route.params.id}`, {
        method: 'PUT',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          hidden: newVal,
        }),
      });
      const content = await resp.json();
      console.log(content);
    },
  },
  mounted: function() {
    /** @type {WebSocket} */
    var socket = new WebSocket(WS_BASE + this.$route.params.id);
    socket.onmessage = (event) => {
      console.log(event);
    }
    socket.onopen = (event) => {
      console.log("Open", event);

      this.fetchBoard();

      /* socket.send(JSON.stringify({
        action: "post",
        data: {
          type: "text",
          string: "Hello",
          x: 100,
          y: 50,
        },
      })); */
    };
    this.socket = socket;
  },
  unmounted: function() {
      /** @type {WebSocket} */
      var socket = this.socket;
      socket.close();
  },
  methods: {
    fetchBoard: function () {
      fetch(`${API_BASE}/boards/${this.$route.params.id}`)
        .then(resp => {
          if (resp.ok) {
            return resp.json()
          }
        })
        .then(json => {
          if (json.data.content == null) {
            json.data.content = [];
          }

          this.board = json.data;
          this.hidden = this.board.hidden;
        });
    },
    deleteBoard: async function() {
      const resp = await fetch(`${API_BASE}/boards/${this.$route.params.id}`, {
        method: 'DELETE',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
      });
      const content = await resp.json();
      if (content.data == 'OK') {
        this.$router.push({name: "Home"});
      }
    },
    addText: function() {
      const textItem = {
        type: "text",
        string: "Hello",
        x: 100,
        y: 50,
      };
      this.board.content.push(textItem);
      this.sendUpdate('post', textItem);
      /*
      */
    },
    sendUpdate: async function(action, data) {
        const resp = await fetch(`${API_BASE}/boards/${this.$route.params.id}`, {
          method: 'POST',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            action,
            data,
          }),
        });
        const content = await resp.json();

        console.log(content);
    }
  },
}
</script>
