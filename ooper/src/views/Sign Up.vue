<template>
  <h1>new to ooper?</h1>
    <div id="toggle"><Toggle @updatePassenger="updatePassenger"/></div>
  <div class="row">
      <input type="text" placeholder="first name" required/>
      <input type="text" placeholder="last name" required/>
  </div>
  <div class="row">
      <input type="text" placeholder="mobile number" required/>
      <input type="email" placeholder="e-mail address" required/>
  </div>
  <div class="row" id="driver-only" :class="{hide:isPassenger,show:isPassenger==false}">
      <input type="text" placeholder="identification number" required :disabled="isPassenger"/>
      <input type="email" placeholder="car license number" required :disabled="isPassenger"/>
  </div>
  <Button text="sign up"/>
</template>

<script>
import Toggle from "../components/toggle.vue"
import Button from "../components/button.vue"
export default {
    components:{Toggle,Button},
    data(){
        return{
            isPassenger:true
        }
    },
    methods:{
        updatePassenger(passenger){
            this.isPassenger = passenger
        }
    },
}
</script>

<style scoped>
h1{
    font-size:54px;
    color:var(--dark-yellow);
    font-weight:normal;
}
.row{
    display:flex;
    justify-content: center;
}
input{
    margin:0 30px 30px 30px;
}
#toggle{
    display: flex;
    justify-content: center;
    margin-top:100px;
}
::v-deep(button){
    margin-top:40px;
}
#driver-only{
    transition: height 0.1s;
}
#driver-only.show{
    height:74px;
}
#driver-only.hide{
    height:0px;
}
</style>