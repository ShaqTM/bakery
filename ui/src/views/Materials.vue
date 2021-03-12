<template>
      <v-card
    class="mx-auto"
    max-width="1000"
    tile
  >

      <v-data-table
        :headers="headers"
        :items="items"
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
        <Material v-bind:content="mdata"/>
        <MaterialPrice v-bind:content="mdata"/>
    </v-card> 
</template>
 
<script>
import Material from "../components/Material"
import MaterialPrice from "../components/MaterialPrice"
import { mdiCurrencyUsd } from '@mdi/js';
  export default {
    name: 'Materials',
    components: {Material,MaterialPrice},
    computed:{
      items(){
        return this.$store.getters.getMaterials;
      }
    },
    data(){
      var mdata = {
      dialog:false,
      editPrice:false,
      id:-1,
      name:"",
      recipe_unit_id:-1,
      recipe_unit_name:"",
      recipe_unit_short_name:"",
      pricet_id:-1,
      price_name:"",
      price_unit_short_name:"",
      coefficient:1}
      var headers= [
          {
            text: 'Наименование',
            align: 'start',
            sortable: true,
            value: 'name',
          },
          { text: 'Ед. рецепта', value: 'recipe_unit_short_name' },
          { text: 'Ед. цены', value: 'price_unit_short_name' },
          { text: 'Коэф. пересчета', value: 'coefficient' },
          { text: 'Цена', value: 'price' },
          { text: '--------', value: 'actions', sortable: false },

        ]
      var icons= {
      mdiCurrencyUsd
      }
      return {mdata:mdata,headers:headers,icons:icons}
    },
    created() {
      this.$store.dispatch('readMaterials',true)
    },
    methods:{
      openDialog(id){
        if (id==-1) {
          this.mdata = {
            dialog:true,
            editPrice:false,
            id:-1,
            name:"",
            recipe_unit_id:-1,
            recipe_unit_name:"",
            recipe_unit_short_name:"",
            price_unit_id:-1,
            price_unit_name:"",
            price_unit_short_name:"",
            coefficient:1}   
        }else{
          this.$store.dispatch('readMaterial', {id:id,price:false})
            .then(resp=>{
              resp.data['dialog']=true        
              this.mdata = resp.data
            })
            .catch(err => console.log(err))
        }
      },
      editPrice(id){
        this.$store.dispatch('readMaterial', {id:id,price:true})
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
