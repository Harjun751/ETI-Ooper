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
const Swal = require('sweetalert2')
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
            var data = {"email":this.email,"password":this.password,"isPassenger":this.isPassenger}
            await fetch(process.env.VUE_APP_AUTH_MS_HOST+"/api/v1/login",{
                body: JSON.stringify(data),
                method:"POST",
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials:'include',
            })
            .then(async (res)=> {
                if (res.status==403){
                    Swal.fire({
                        title: 'failed...',
                        text: 'incorrect username or password',
                        icon: 'error',
                        confirmButtonText: 'close',
                        customClass:{
                            popup: 'custom-swal-modal',
                            icon: 'custom-swal-icon',
                            content: 'custom-swal-content',
                            confirmButton: 'custom-swal-button'
                        }
                    })
                }
                return
            })
            .then(()=>{
                // store.setJWTAccessToken(data.token)
                store.setIsPassenger(this.isPassenger)
            })
            .then(()=>{
                if (this.isPassenger){
                this.$router.push("new-trip")
                }
                else if (!this.isPassenger){
                    this.$router.push("trip-management")
                }
            })
            .catch(()=>{
                Swal.fire({
                    title: 'failed...',
                    text: 'failed to login due to server issues',
                    icon: 'error',
                    confirmButtonText: 'close',
                    customClass:{
                        popup: 'custom-swal-modal',
                        icon: 'custom-swal-icon',
                        content: 'custom-swal-content',
                        confirmButton: 'custom-swal-button'
                    }
                })
            })
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