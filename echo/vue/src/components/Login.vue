<template>
  <div class="app">
    <h1>Login</h1>
    <div class="form">
      <input type="text" id="email" placeholder="Email"/>
      <input type="password" id="password" placeholder="Password"/>
      <button type="button" v-on:click="login">Login</button>
    </div>
    <h5>{{ error }}</h5>
  </div>
</template>

<script>
import { server, post, validateEmail } from '../res'
export default {
  name: 'login',
  data () {
    return {
      error: null,
    }
  },
  methods: {
    login: function () {
      const email = document.getElementById('email').value;
      const password = document.getElementById('password').value;

      if( validateEmail(email) && password.length > 5 ) {
        this.error = null;

        post(
          server + 'login',
          ['email', 'password'],
          [email, password],
        )
        .then(res => {
          console.log(res);
          if(res.done) {
            this.$router.push('/');
          } else {
            this.error = res.msg
          }
        })
        .catch(res => this.error = res)
      } else {
        this.error = 'An error occured'
      }
    },
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.app {
  display: flex;
  flex-flow: column;
  justify-content: center;
  align-items: center;
}

.form {
  display: flex;
  flex-flow: column;
  align-items: center;
}

.form * {
  margin-top: 20px;
  width: 90vw;
  max-width: 300px;
  height: 26px;
}

.form button {
  border: none;
  border-radius: 3px;
  background: #317dba;
  color: #fff;
}

</style>
