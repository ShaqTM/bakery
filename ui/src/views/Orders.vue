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
        <Order v-bind:content="mdata"/>
    </v-card> 
</template>
 
<script>
import Order from "../components/Order"
import { mdiCurrencyUsd } from '@mdi/js';
  export default {
    name: 'Orders',
    components: {OrderOrder},
    computed:{
      items(){
        return this.$store.getters.getOrders;
      }
    },
    data(){
      var mdata = {
      dialog:false,
      id:-1,
      name:"",
      unit_id:-1,
      unit_short_name:"",
      output:0,
      content:[]}
      var headers= [
          {
            text: 'Заказчик',
            align: 'start',
            sortable: true,
            value: 'customer',
          },
          { text: 'Рецепт', value: 'recipe_name'},
          { text: 'Дата', value: 'date'},
          { text: 'Дата выдачи', value: 'release_date'},
          { text: 'Ед. изм', value: 'unit_name' },
          { text: 'Количество план', value: 'plan_qty'},
          { text: 'Цена', value: 'price'},
          { text: 'Стоимость план', value: 'plan_cost'},
          { text: 'Количество факт', value: 'fact_qty'},
          { text: 'Стоимость факт', value: 'fact_cost'},
          { text: 'Стоимость материалов', value: 'materials_cost'},
          { text: '--------', value: 'actions', sortable: false },
 
        ]
      var icons= {
      mdiCurrencyUsd
      }
      return {mdata:mdata,headers:headers,icons:icons}
    },
    created() {
      this.$store.dispatch('readOrders')
    },
    methods:{
      openDialog(id){
        if (id==-1) {
          this.mdata = {
            dialog:true,
            id:-1,
            name:"",
            unit_id:-1,
            unit_short_name:"",
            output:0,
            content:[]}   
        }else{
          this.$store.dispatch('readOrder', {id:id})
            .then(resp=>{
              resp.data['dialog']=true        
              this.mdata = resp.data
            })
            .catch(err => console.log(err))
        }
      },
      editPrice(id){
        this.$store.dispatch('readOrder', {id:id})
          .then(resp=>{
            resp.data['dialog']=false
            this.mdata = resp.data
            })
          .catch(err => console.log(err))
      }
           
    }
  }
</script>
