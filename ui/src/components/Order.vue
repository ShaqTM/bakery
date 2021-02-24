<template>
    <v-dialog
      v-model="this.content.dialog"
      persistent
      max-width="1000"
    >
       <v-card>
        <v-card-title>
          <span class="headline">Заказ</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Номер заказа"
                  required
                  readonly
                  v-model="content.id"
                ></v-text-field>
              </v-col>
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Дата"
                  required
                  v-model="content.date"
                ></v-text-field>
              </v-col>
           </v-row>                       
            <v-row>
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-select
                  v-model="content.recipe_id"
                  :items="recipes"
                  item-text="name"
                  item-value="id"
                  label="Рецепт"
                  required
                  @change="recipeChanged()"
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
                  label="Заказчик"
                  required
                  v-model="content.customer"
                ></v-text-field>
              </v-col>
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Дата выдачи"
                  required
                  v-model="content.customer"
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
                  label="Цена"
                  required
                  v-model.number="content.price"
                  :rules="[rules.num]"
                  @change="priceChanged()"
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
      }

    },
    created() {
      this.$store.dispatch('readUnits')
      this.$store.dispatch('readMaterials',true)
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
      return {rules:rules,headers:headers,icons:icons}
    },
//    data(){
      
//        if (this.content.id==-1) {
//          return {short_name:"",name:""};
//        }

//        return  this.$store.getters.getUnit(this.content.id);
//    },
    methods:{
      saveData(){
        if (isNaN(this.content.output)){
          return
        }        
        var len = this.content.content.length
        for(let i = 0 ; i < len; i++) {
          this.content.content[i]["string_order"] = i
        }        
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
        item.price = mMaterial.price
        item.cost = item/price*item.qty
        countMaterialCost()

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
        for(let i = 0 ; i < len; i++) {
          this.content.content[i]["string_order"] = i
        }        
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
        for(let i = 0 ; i < len; i++) {
          this.content.content[i]["string_order"] = i
        }
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
        len = this.content.content.length
        for(let i = 0 ; i < len; i++) {
          this.content.content[i]["string_order"] = i
        }
      },
      countMaterialCost(){
        var mMaterialCost = 0
        var len = this.content.content.length
        for(let i = 0 ; i < len; i++) {
          mMaterialCost += this.content.content[i]["cost"]
        } 
        this.content.materials_cost = mMaterialCost       
      },
      materialQtyCHanged(item){
        item.cost=item.qty*item.price
        countMaterialCost()
      },
      materialPriceChanged(item){
        item.cost=item.qty*item.price
        countMaterialCost()
      },
      recipeChanged(){
        var len = this.content.content.length
        for(let i = 0 ; i < len; ) {
          if (this.content.content[i]["by_recipe"]===true){
            this.content.content.splice(mIndex, 1)
            len--
          }else{
            i++
          }
        }
        this.$store.dispatch('readRecipe', {id:id,price:true})
          .then(resp=>{
            var len = resp.data.content.length
            for(let i = 0 ; i < len; i++) {
              this.content.content.splice(i, 0,resp.data.content[i])
            }
          })         
          .catch(err => console.log(err))
      }
      @change="priceChanged(item)"
      @change="planQtyChanged(item)"
      @change="planCostChanged(item)"
      @change="factCostChanged(item)"
      @change="factQtyChanged(item)"
    }
  }
</script>
