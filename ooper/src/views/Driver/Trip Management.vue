<template>
    <div v-if="exists">
        <div id="map">
            <iframe
            width="680"
            height="600"
            frameborder="0" style="border:3px solid var(--bright-yellow);border-radius:50px;"
            :src="fullURL" allowfullscreen>
            </iframe>
        </div>
        <div id="form">
            <input v-model="origin" type="text" placeholder="start point" disabled/>
            <input v-model="destination" type="text" placeholder="end point" disabled/>
            <span v-if="start!=false">
                <Button text="start trip" disabled="true" @click="startTrip"/>
            </span>
            <span v-else>
                <Button text="start trip" @click="startTrip"/>
            </span>
            <br/>
            <span v-if="start==false || end!=false">
                <Button text="end trip" disabled="true" @click="endTrip"/>
            </span>
            <span v-else>
                <Button text="end trip" @click="endTrip"/>
            </span>
        </div>
    </div>
    <div v-if="exists==false">
        <h3 style="color:var(--bright-yellow);">No Trips!</h3>
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
            start:"",
            end:"",
            tripID:"",
            exists:null,
            fullURL:""
        }
    },
    async mounted(){
        await fetch(process.env.VUE_APP_TRIP_MS_HOST+"/api/v1/current-trip",{
        method:"GET",
        headers: {
            'Content-Type': 'application/json',
            'Authorization': "Bearer " +  store.state.jwtAccessToken
        },
        credentials:'include',
        })
        .then(async (res)=>{
            if (res.status==404){
                this.exists = false
                throw "No data"
            }
            return await res.json()
        })
        .then((data)=>{
            this.exists = true
            this.origin = data.PickUp
            this.destination = data.DropOff
            this.start = data.Start.Valid
            this.end = data.End.Valid
            this.tripID = data.ID
            this.fullURL = "https://www.google.com/maps/embed/v1/directions?key=" + process.env.VUE_APP_GMAPS_KEY + "&origin="+this.origin+"&destination="+this.destination
        })
    },
    methods:{
        async startTrip(){
            await fetch(process.env.VUE_APP_TRIP_MS_HOST+"/api/v1/trip/"+this.tripID+"/start",{
                method:"POST",
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': "Bearer " +  store.state.jwtAccessToken
                },
                credentials:'include',
            })
            .then(()=>{
                var today = new Date();
                this.start = today.getDate() + today.getTime()
            })
        },
        async endTrip(){
            await fetch(process.env.VUE_APP_TRIP_MS_HOST+"/api/v1/trip/"+this.tripID+"/end",{
                method:"POST",
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': "Bearer " +  store.state.jwtAccessToken
                },
                credentials:'include',
            })
            .then(()=>{
                var today = new Date();
                this.end = today.getDate() + today.getTime()
                this.exists = false
                Swal.fire({
                    title: 'finished trip!',
                    text: "you finished this trip!",
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
        },
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