<template>
  <div class="app">
    <h1>Sign up</h1>
    <div class="form">
      <input type="text" id="name" placeholder="Name" autofocus=""/>
      <input type="text" id="email" placeholder="Email"/>
      <input type="password" id="password" placeholder="Password"/>
      <input type="password" id="password2" placeholder="Repeat Password"/>
      <button type="button" v-on:click="signup">Sign up</button>
    </div>
    <h5>{{ error }}</h5>
  </div>
</template>

<script>
import { server, post, validateEmail } from '../res'
export default {
  name: 'signup',
  data () {
    return {
      error: null,
    }
  },
  methods: {
    signup: function () {
      const name = document.getElementById('name').value;
      const email = document.getElementById('email').value;
      const password = document.getElementById('password').value;
      const password2 = document.getElementById('password2').value;

      if( validateEmail(email) && name.length > 2 && password.length > 5 && password === password2 ) {
        this.error = null;

        post(
          server + 'signup',
          ['name', 'email', 'password'],
          [name, email, password],
        )
        .then(res => {
          if(res.token) {
            localStorage.setItem('token', res.token);
            // this.$router.push('/');
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
