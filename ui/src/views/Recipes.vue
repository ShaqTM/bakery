<template>
      <v-card
    class="mx-auto"
    max-width="1000"
    tile
  >

      <v-data-table
        :headers="headers"
        :items="items"
        :items-per-page="15"
        class="elevation-1"
        disable-pagination=true
        hide-default-footer=true        
      >
        <template v-slot:[`item.actions`]="{ item }">
          <v-icon
            small
            class="mr-2"
            @click="openDialog(item.id)"
          >
            mdi-pencil
          </v-icon>
          <v-icon
            small
            class="mr-2"
            @click="editPrice(item.id)"
          >
            {{ icons.mdiCurrencyUsd }}
          </v-icon>          
        </template>      


      </v-data-table>
      <div class="text-center pt-2">
        <v-btn
          color="primary"
          @click.stop="openDialog(-1)"
        >
          +
        </v-btn>
      </div>
        <Recipe v-bind:content="mdata"/>
        <RecipePrice v-bind:content="mdata"/>
    </v-card> 
</template>
 
<script>
import Recipe from "../components/Recipe"
import RecipePrice from "../components/RecipePrice"
import { mdiCurrencyUsd } from '@mdi/js';
  export default {
    name: 'Recipes',
    components: {Recipe,RecipePrice},
    computed:{
      items(){
        return this.$store.getters.getRecipes;
      }
    },
    data(){
      var mdata = {
      dialog:false,
      editPrice:false,
      id:-1,
      name:"",
      unit_id:-1,
      unit_short_name:"",
      output:0,
      content:[]}
      var headers= [
          {
            text: 'Наименование',
            align: 'start',
            sortable: true,
            value: 'name',
          },
          { text: 'Выход', value: 'output' },
          { text: 'Ед. изм', value: 'unit_short_name' },
          { text: 'Цена', value: 'price' },
          { text: '--------', value: 'actions', sortable: false },

        ]
      var icons= {
      mdiCurrencyUsd
      }
      return {mdata:mdata,headers:headers,icons:icons}
    },
    created() {
      this.$store.dispatch('readRecipes',true)
    },
    methods:{
      openDialog(id){
        if (id==-1) {
          this.mdata = {
            dialog:true,
            editPrice:false,
            id:-1,
            name:"",
            unit_id:-1,
            unit_short_name:"",
            output:0,
            content:[]}   
        }else{
          this.$store.dispatch('readRecipe', {id:id,price:false})
            .then(resp=>{
              resp.data['dialog']=true        
              this.mdata = resp.data
            })
            .catch(err => console.log(err))
        }
      },
      editPrice(id){
        this.$store.dispatch('readRecipe', {id:id,price:true})
          .then(resp=>{
            resp.data['dialog']=false
            resp.data['editPrice']=true
            this.mdata = resp.data
            })
          .catch(err => console.log(err))
      }
           
    }
  }
</script>
