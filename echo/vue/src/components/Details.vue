<template>
  <div>
    <div v-if="item" class="item">
      <h2>{{item.title}}</h2>
      <img v-bind:src="item.src"/>
      <p>{{item.desc}}</p>
      <hr />
      <div v-if="item.comments.length > 0">
        <p v-for="c in item.comments">
          {{c}}
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import {server} from '../res';
export default {
  name: 'details',
  data () {
    return {
      msg: 'Details',
      item: null,
    }
  },
  beforeMount(){
    if(typeof this.$router.history.current.query.data !== 'string') {
      this.item = this.$router.history.current.query.data;
    } else {
      const addr = this.$router.history.current.path
      const result = addr.substring(addr.lastIndexOf("/") + 1);
      // fetch('https://react-eko.herokuapp.com/api/' + result)
      // .then(res => res.json())
      // .then(res => {
      //   this.item = res
      // })
    }
 },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.item {
  display: flex;
  flex-flow: column;
  justify-content: center;
  align-items: center;
}

.item * {
  width: 90vw;
  max-width: 500px;
  margin-top: 20px;
}

</style>
