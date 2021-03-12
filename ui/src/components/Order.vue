<template>
    <v-dialog
      v-model="this.content.dialog"
      persistent
      max-width="1200"
    >
       <v-card>
        <v-card-title>
          <span class="headline">Заказ</span>
        </v-card-title>
        <v-card-text>
          <v-container dense>
            <v-row>
              <v-col

              >
                <v-text-field
                  label="Номер заказа"
                  required
                  readonly
                  v-model="content.id"
                  dense
                ></v-text-field>
              </v-col>
              <v-col

              >
                <v-menu
                  v-model="menu1"
                  :close-on-content-click="false"
                  :nudge-right="40"
                  transition="scale-transition"
                  offset-y
                  min-width="290px"
                >
                  <template v-slot:activator="{ on, attrs }">
                    <v-text-field
                      v-model="formatedDate"
                      label="Дата"
                      readonly
                      v-bind="attrs"
                      v-on="on"
                      dense
                      required
                    ></v-text-field>
                  </template>
                  <v-date-picker v-model="content.date" locale="ru-RU" @input="menu1 = false"></v-date-picker>
                </v-menu>
              </v-col>
           </v-row>                       
              <v-select
                  v-model="content.recipe_id"
                  :items="recipes"
                  item-text="name"
                  item-value="id"
                  label="Рецепт"
                  required
                  @change="recipeChanged()"
                  dense
                ></v-select>
            <v-row>
              <v-col
              >
                <v-text-field
                  label="Заказчик"
                  required
                  v-model="content.customer"
                  dense
                ></v-text-field>
              </v-col>
              <v-col
                cols="12"
                sm="3"
                md="2"
              >
               <v-menu
                  v-model="menu2"
                  :close-on-content-click="false"
                  :nudge-right="40"
                  transition="scale-transition"
                  offset-y
                  min-width="290px"
                >
                  <template v-slot:activator="{ on, attrs }">
                    <v-text-field
                      v-model="formatedReleaseDate"
                      label="Дата выдачи"
                      readonly
                      v-bind="attrs"
                      v-on="on"
                      dense
                      required
                    ></v-text-field>
                  </template>
                  <v-date-picker v-model="content.release_date"  locale="ru-RU" @input="menu2 = false"></v-date-picker>
                </v-menu>
              </v-col>              
           </v-row>           
           <v-row>
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Цена"
                  required
                  v-model.number="content.price"
                  :rules="[rules.num]"
                  @change="priceChanged()"
                  dense
                ></v-text-field>
              </v-col>             
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-select
                  v-model="content.unit_id"
                  :items="units"
                  item-text="short_name"
                  item-value="id"
                  label="Ед. изм."
                  required
                  readonly
                  dense
                ></v-select>
              </v-col>
           </v-row>
           <v-row>
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Плановое количество"
                  required
                  v-model.number="content.plan_qty"
                  :rules="[rules.num]"
                  @change="planQtyChanged()"
                  dense
                ></v-text-field>
              </v-col>             
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Плановая стоимость"
                  required
                  v-model.number="content.plan_cost"
                  :rules="[rules.num]"
                  @change="planCostChanged()"
                  dense
                ></v-text-field>
              </v-col>
           </v-row>   
           <v-row>
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Фактическое количество"
                  required
                  v-model.number="content.fact_qty"
                  :rules="[rules.num]"
                  @change="factQtyChanged()"
                  dense
                ></v-text-field>
              </v-col>             
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Фактическая стоимость"
                  required
                  v-model.number="content.fact_cost"
                  :rules="[rules.num]"
                  @change="factCostChanged()"
                  dense
                ></v-text-field>
              </v-col>
           </v-row>            
           <v-row>
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Стоимость материалов"
                  required
                  v-model.number="content.materials_cost"
                  :rules="[rules.num]"
                  readonly
                  dense
                ></v-text-field>
              </v-col>             
           </v-row>  
          </v-container>
      </v-card-text>
  
     <v-data-table
        dense
        :headers="headers"
        :items="content.content"
        class="elevation-1"
                disable-pagination=true
        hide-default-footer=true
      >
        <template v-slot:[`item.material_id`]="{ item }">
          <v-select
            dense
            v-model="item.material_id"
            :items="materials"
            item-text="name"
            item-value="id"
            required
            @change="materialSelected(item)"
            ></v-select>          
        </template>      
        <template v-slot:[`item.qty`]="{ item }">
         <v-text-field
          dense
            v-model.number="item.qty"
            :rules="[rules.num]"            
            required
            @change="materialQtyChanged(item)"
            ></v-text-field>          
        </template>      
        <template v-slot:[`item.price`]="{ item }">
         <v-text-field
          dense
            v-model.number="item.price"
            :rules="[rules.num]"            
            required
            @change="materialPriceChanged(item)"
            ></v-text-field>          
        </template> 
        <template v-slot:[`item.actions`]="{ item }">
          <v-icon
            small
            class="mr-2"
            @click="moveUp(item)"
          >
            {{ icons.mdiArrowUpThick  }}
          </v-icon>
          <v-icon
            small
            class="mr-2"
            @click="moveDown(item)"
          >
            {{ icons.mdiArrowDownThick  }}
          </v-icon> 
          <v-icon
            small
            class="mr-2"
            @click="deleteRow(item)"
          >
            {{ icons.mdiCloseThick   }}
          </v-icon>                    
        </template>      


      </v-data-table>
      <div class="text-center pt-2">
        <v-btn
          color="primary"
          @click="AddMaterial()"
        >
          +
        </v-btn>
      </div>

          <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="blue darken-1"
            text
            @click="close()"
          >
            Закрыть
          </v-btn>
          <v-btn
            color="blue darken-1"
            text
            @click="saveData()"
          >
            Сохранить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
</template>
 
<script>
  import { mdiArrowDownThick } from '@mdi/js';
  import { mdiArrowUpThick } from '@mdi/js';
  import { mdiCloseThick  } from '@mdi/js';
  export default {
    name: 'Order',
    props:{content:Object
    },
    computed:{
      units(){
        return this.$store.getters.getUnits;
      },
      materials(){
        return this.$store.getters.getMaterials;
      },
      recipes(){
        return this.$store.getters.getRecipes;
      },
      formatedDate(){
        return this.formatDate(this.content.date)
      },
      formatedReleaseDate(){
        return this.formatDate(this.content.release_date)
      }
    },
    created() {
      this.$store.dispatch('readUnits')
      this.$store.dispatch('readMaterials',true)
      this.$store.dispatch('readRecipes',true)
    },    
    data(){
      var headers= [

      {
        text: "Материал",
        align: 'start',
        sortable: true,
        value: 'material_id',
      },
      { text: 'Количество', value: 'qty' },
      { text: 'Ед. изм', value: 'unit_short_name' },
      { text: 'Цена', value: 'price' },
      { text: 'Стоимость', value: 'cost'},
      { text: '--------', value: 'actions', sortable: false },

        ]
      var rules = {num: value => {return !isNaN(value)||'Должно быть число'}}
      var icons= {
      mdiArrowDownThick,
      mdiArrowUpThick,
      mdiCloseThick 

      }
     
      return {rules:rules,headers:headers,icons:icons,menu1: false,menu2: false,recipe:null}
    },
//    data(){
      
//        if (this.content.id==-1) {
//          return {short_name:"",name:""};
//        }

//        return  this.$store.getters.getUnit(this.content.id);
//    },
    methods:{
      formatDate (date) {
        if (!date) return null

        const [year, month, day] = date.split('-')
        return `${day}.${month}.${year}`
      },      
      updateTableOrder(){
        var len = this.content.content.length        
        for(let i = 0 ; i < len; i++) {
          this.content.content[i]["string_order"] = i
        }         
      },
      saveData(){
//        if (isNaN(this.content.output)){
//          return
//        }        
        this.updateTableOrder()       
        this.$store.dispatch('writeOrder',this.content)
        this.content.dialog= false

      },
      close(){
        this.content.dialog = false
      },

      AddMaterial(){
        this.content.content.push({id:this.content.id,material_id:-1,qty:0,price:0,cost:0,unit_id:-1,string_order:this.content.content.length,by_recipe:false})
      },
      materialSelected(item){
        var mMaterial=this.$store.getters.getMaterial(item.material_id)
        item.unit_id = mMaterial.recipe_unit_id
        var mUnit = this.$store.getters.getUnit(item.unit_id)
        item.unit_short_name = mUnit.short_name
        if (mMaterial.coefficient==0){
          item.price = 0
        } else{
          item.price = Math.round(mMaterial.price/mMaterial.coefficient*10000)/10000
        }

        
        item.cost = item.price*item.qty
        this.countMaterialCost()

      },
      moveUp(item){
        var len = this.content.content.length
        var mIndex
        for(let i = 0 ; i < len; i++) {
          if (this.content.content[i]["string_order"] === item.string_order) {
            mIndex = i
            break
          }
        }
        if (mIndex===0){
          return
        }

        this.content.content.splice(mIndex, 1, this.content.content[mIndex-1])
        this.content.content.splice([mIndex-1],1,item)
        this.updateTableOrder()
       
      },
      moveDown(item){
        var len = this.content.content.length
        var mIndex
        for(let i = 0 ; i < len; i++) {
          if (this.content.content[i]["string_order"] === item.string_order) {
            mIndex = i
            break
          }
        }
        if (mIndex===len-1){
          return
        }
        this.content.content.splice(mIndex, 1, this.content.content[mIndex+1])
        this.content.content.splice([mIndex+1],1,item)
        this.updateTableOrder()
      } ,
      deleteRow(item){
        var len = this.content.content.length
        var mIndex
        for(let i = 0 ; i < len; i++) {
          if (this.content.content[i]["string_order"] === item.string_order) {
            mIndex = i
            break
          }
        }
        this.content.content.splice(mIndex,1)
        this.updateTableOrder()
      },
      countMaterialCost(){
        var mMaterialCost = 0
        var len = this.content.content.length
        for(let i = 0 ; i < len; i++) {
          mMaterialCost += this.content.content[i]["cost"]
        } 
        this.content.materials_cost =  Math.round(mMaterialCost*100)/100
      },
      materialQtyChanged(item){
        item.cost=item.qty*item.price
        this.countMaterialCost()
      },
      materialPriceChanged(item){
        item.cost=item.qty*item.price
        this.countMaterialCost()
      },
      recipeChanged(){

        this.$store.dispatch('readRecipe', {id:this.content.recipe_id,price:true})
          .then(resp=>{
            this.recipe = resp.data
            this.updateOrderByRecipe()
          })         
          .catch(err => console.log(err))
      },
      updateOrderByRecipe(){
        var len = this.content.content.length
        for(let i = 0 ; i < len; ) {
          if (this.content.content[i]["by_recipe"]===true){
            this.content.content.splice(i, 1)
            len--
          }else{
            i++
          }
        }        
        len = this.recipe.content.length
        for(let i = 0 ; i < len; i++) {
            this.recipe.content[i]["by_recipe"] = true
            this.content.content.splice(i, 0,Object.assign({}, this.recipe.content[i]))
        }
        len = this.content.content.length
        for(let i = 0 ; i < len; i++) {
          if (this.content.content[i]["by_recipe"]){
            this.content.content[i].qty = Math.round(this.recipe.content[i].qty*this.content.plan_qty/this.recipe.output*100)/100
            this.content.content[i].cost = Math.round(this.content.content[i].price*this.content.content[i].qty*100)/100
          }
        }
        this.updateTableOrder()
        this.countMaterialCost()
        this.content.price = this.recipe.price
        this.content.plan_cost =  Math.round(this.content.price*this.content.plan_qty*100)/100
        this.content.unit_short_name = this.recipe.unit_short_name
        this.content.unit_id = this.recipe.unit_id        

      },



      priceChanged(){
        this.content.plan_cost = this.content.price*this.content.plan_qty
        this.content.fact_cost = this.content.price*this.content.fact_qty
      },
      planQtyChanged(){
        if (this.recipe==null){
          this.recipeChanged()          
        } else{
          this.updateOrderByRecipe()
        }
      },
      planCostChanged(){
        if (this.content.plan_qty!=0){
          this.content.price = Math.round(this.content.plan_cost/this.content.plan_qty*100)/100
          this.content.fact_cost = this.content.price*this.content.fact_qty
        }
      },
      factQtyChanged(){
        this.content.fact_cost = this.content.price*this.content.fact_qty
      },
      factCostChanged(){
        if (this.content.fact_qty!=0){
          this.content.price = Math.round(this.content.fact_cost/this.content.fact_qty*100)/100
          this.content.plan_cost = this.content.price*this.content.plan_qty
        }
      }
    }
  }
</script>
