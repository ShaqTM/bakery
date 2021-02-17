<template>
    <v-dialog
      v-model="this.content.dialog"
      persistent
      max-width="1000"
    >
       <v-card>
        <v-card-title>
          <span class="headline">Рецепт</span>
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
                  label="Наименование"
                  required
                  v-model="content.name"
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
                  label="Выход"
                  required
                  v-model.number="content.output"
                  :rules="[rules.num]"
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
                ></v-select>
              </v-col>
           </v-row>
           
          </v-container>
        </v-card-text>
     <v-data-table
        :headers="headers"
        :items="content.content"
        class="elevation-1"
      >
        <template v-slot:[`item.material_id`]="{ item }">
          <v-select
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
            v-model.number="item.qty"
            :rules="[rules.num]"            
            required
            ></v-text-field>          
        </template>      


      </v-data-table>

        <v-fab-transition>
          <v-btn
            color="pink"
            dark
            absolute
            bottom
            right
            fab
            @click.stop="AddMaterial()"
          >
            <v-icon>mdi-plus</v-icon>
          </v-btn>
          
        </v-fab-transition>        
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
  export default {
    name: 'Recipe',
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
      this.$store.dispatch('readMaterials',false)
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
//      { text: '--------', value: 'actions', sortable: false },

        ]
      var rules = {num: value => {return !isNaN(value)||'Должно быть число'}}
      return {rules:rules,headers:headers}
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
        this.$store.dispatch('writeRecipe',this.content)
        this.content.dialog= false

      },
      close(){
        this.content.dialog = false
      },
      AddMaterial(){
        this.content.content.push({id:this.content.id,material_id:-1,qty:0,unit_id:-1,string_order:this.content.content.length})
      },
      materialSelected(item){
        var mMaterial=this.$store.getters.getMaterial(item.material_id)
        item.unit_id = mMaterial.recipe_unit_id
        var mUnit = this.$store.getters.getUnit(item.unit_id)
        item.unit_short_name = mUnit.short_name

      }
    }
  }
</script>
