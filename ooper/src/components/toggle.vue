<template>
    <div>
        <div id="toggler" @click="toggleMode">
            <span :class="{ active: passenger==false}">driver</span>
            <span :class="{ active: passenger==true}">passenger</span>
        </div>
        <div id="slider" :class="{passenger:passenger==true, driver:passenger==false}"></div>
    </div>
</template>

<script>
export default{
    data(){
        return{
            passenger:true
        }
    },
    methods:{
        toggleMode(){
            this.passenger = !this.passenger
            this.$emit("updatePassenger",this.passenger)
        }
    }
}
</script>

<style scoped>
#toggler{
    cursor: pointer;
    border:3px solid var(--bright-yellow);
    border-radius: 50px;
    width:400px;
    height:66px;
    display:flex;
    justify-content: space-between;
    align-items: center;
}
#toggler span{
    -webkit-touch-callout: none;
    -webkit-user-select: none;
    -khtml-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
    margin:0 35px 0 35px;
    color:var(--bright-yellow);
    transition:color 0.1s;
}
#toggler span.active{
    color:var(--purple);
    font-weight:bold;
}
#slider{
    position:relative;
    background:var(--bright-yellow);
    top:-72px;
    height:66px;
    z-index: -1;
    border-radius: 50px;
    border:3px solid var(--bright-yellow);
}
#slider.passenger{
    left:161px;
    width:237px;
    transition: width 0.1s, left 0.1s;
}
#slider.driver{
    left:0px;
    width:181px;
    transition: width 0.1s, left 0.1s;
}
</style>