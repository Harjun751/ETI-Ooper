<template>
    <div id="images">
        <Driver :isPassenger="isPassenger"/>
        <Passenger :isPassenger="isPassenger"/>
    </div>
    <div id="login">
        <Toggle @updatePassenger="updatePassenger"/>
        <input v-model="email" type="email" placeholder="e-mail address" required/>
        <input v-model="password" type="password" placeholder="password" required/>
        <Button text="login" @click="login"/>
    </div>
</template>

<script>
import Toggle from "../components/toggle.vue"
import Button from "../components/button.vue"
import Passenger from "../components/passenger.vue"
import Driver from "../components/driver.vue"
import { store } from "../state"
export default {
    components:{Toggle,Button,Passenger,Driver},
    data(){
        return{
            isPassenger:true,
            email:null,
            password:null
        }
    },
    methods:{
        updatePassenger(passenger){
            this.isPassenger = passenger
        },
        async login(){
            if (this.isPassenger){
                var data = {"email":this.email,"password":this.password}
                await fetch("http://localhost:5000/api/v1/login",{
                    body: JSON.stringify(data),
                    method:"POST",
                    headers: {
                        'Content-Type': 'application/json',
                    },
                })
                .then(async (res)=> await res.json())
                .then((data)=>{
                    store.setJWTAccessToken(data.token)
                    if (data.isPassenger == "true"){
                        store.setIsPassenger(true)
                    }
                    else if (data.isPassenger == "false"){
                        store.setIsPassenger(false)
                    }
                })
                this.$router.push("new-trip")
            }
        }
    },
}
</script>

<style scoped>
#login{
    position:relative;
    float:right;
    top:200px;
    right:50px;
}
input{
    margin:0 auto 30px auto;
}
input:last-of-type{
    margin-bottom: 60px;
}
#images{
    float:left;
    margin-top:200px;
    margin-left:80px;
}
</style>