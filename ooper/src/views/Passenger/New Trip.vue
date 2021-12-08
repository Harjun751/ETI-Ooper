<template>
  <div id="map">
    <div v-if="mapskey!=nil">
        <iframe
        frameborder="0"
        :src="fullURL" allowfullscreen>
        </iframe>
    </div>
  </div>
  <div id="form" :class="{ center: mapskey==nil }">
    <input v-model="origin" type="text" placeholder="start point"/>
    <input v-model="destination" type="text" placeholder="end point"/>
    <p>fee:<span id="value">{{ price }}</span></p>
    <Button text="book" @click="requestTrip"/>
  </div>
</template>

<script>
import Button from "../../components/button.vue"
import { store } from "../../state"
const Swal = require('sweetalert2')
export default {
    components:{Button},
    data(){
        return{
            origin:"",
            destination:"",
            price:"",
            mapskey:process.env.VUE_APP_GMAPS_KEY
        }
    },
    computed:{
        fullURL(){
            if (this.origin != "" && this.destination != ""){
                return "https://www.google.com/maps/embed/v1/directions?key=" + process.env.VUE_APP_GMAPS_KEY + "&origin="+this.origin+"&destination="+this.destination
            }
            else{
                return "https://www.google.com/maps/embed/v1/view?key=" + process.env.VUE_APP_GMAPS_KEY +"&maptype=roadmap&center=1.38229,103.79714&zoom=10"
            }
        }
    },
    watch:{
        destination(){
            if (this.origin != "" && this.destination != ""){
                this.price = "$" + (Math.floor(Math.random()*100)).toString()
            }
        },
        origin(){
            if (this.origin != "" && this.destination != ""){
                this.price = "$" + (Math.floor(Math.random()*100)).toString()
            }
        },
    },
    methods:{
        async requestTrip(){
            var data = {"PickUp":this.origin,"DropOff":this.destination}
            await fetch(process.env.VUE_APP_TRIP_MS_HOST+"/api/v1/trips",{
                body: JSON.stringify(data),
                method:"POST",
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': "Bearer " +  store.state.jwtAccessToken
                },
                credentials:'include',
            })
            .then(async (res)=> {
                if (res.status==404){
                    Swal.fire({
                        title: 'failed...',
                        text: 'no drivers available! try again soon...',
                        icon: 'error',
                        confirmButtonText: 'close',
                        customClass:{
                            popup: 'custom-swal-modal',
                            icon: 'custom-swal-icon',
                            content: 'custom-swal-content',
                            confirmButton: 'custom-swal-button'
                        }
                    })
                    return
                }
                return await res.json()
            })
            .then((data)=>{
                Swal.fire({
                    title: 'done!',
                    text: "Your driver is " + data.FirstName + " " + data.LastName + "\nLicense Number: " + data.LicenseNumber,
                    icon: 'success',
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
    }
}
</script>

<style scoped>
#form{
    float:right;
    margin-top:200px;
    margin-right:30px;
}
#form p{
    color:var(--bright-yellow);
    text-align: left;
    margin-bottom:50px;
}
#form.center{
    margin-left:auto;
    margin-right:auto;
    display:inline-block;
    float:none;
}
#value{
    padding-left:10px;
    text-decoration: underline;
}
#map{
    float:left;
    margin-top:100px;
    margin-left:150px;
}
input{
    width:700px;
}
iframe{
    border:3px solid var(--bright-yellow);
    border-radius:50px;
    width:680px;
    height:600px;
}
@media screen and (max-width: 1440px) {
  #map{
    margin-left:30px;
  }
  iframe{
    width:580px;
    height:500px;
  }
  input{
    width:500px;
  }
}
</style>