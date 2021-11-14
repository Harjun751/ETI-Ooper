<template>
  <div>
    <button type="submit">{{ text }}</button>
  </div>
</template>

<script>
export default {
    props:["text"],
    data(){
        return{
            
        }
    },
    methods:{

    }
}
</script>

<style>
button{
    all:unset;
    cursor: pointer;
    border-radius: 50px;
    width:268px;
    height:68px;
    background:var(--bright-yellow);
    color:var(--purple);
    font-weight: bold;
}
</style>