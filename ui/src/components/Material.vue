<template>
    <v-dialog
      v-model="this.content.dialog"
      persistent
      max-width="500"
    >
       <v-card>
        <v-card-title>
          <span class="headline">Материал</span>
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
                <v-select
                  v-model="content.price_unit_id"
                  :items="units"
                  item-text="short_name"
                  item-value="id"
                  label="Ед. изм. цены"
                  required
                ></v-select>
              </v-col>
             <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-select
                  v-model="content.recipe_unit_id"
                  :items="units"
                  item-text="short_name"
                  item-value="id"
                  label="Ед. изм. рецепта"
                  required
                ></v-select>
              </v-col>
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Коэфф. пересчета"
                  required
                  v-model.number="content.coefficient"
                  :rules="[rules.num]"
                ></v-text-field>
              </v-col>
               
           </v-row>
           
          </v-container>
        </v-card-text>
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
    name: 'Unit',
    props:{content:Object
    },
    computed:{
      units(){
        return this.$store.getters.getUnits;
      }
    },
    created() {
      this.$store.dispatch('readUnits')
    },    
    data(){
      return {rules:{num: value => {return !isNaN(value)||'Должно быть число'}}}
    },
//    data(){
      
//        if (this.content.id==-1) {
//          return {short_name:"",name:""};
//        }

//        return  this.$store.getters.getUnit(this.content.id);
//    },
    methods:{
      saveData(){
        if (isNaN(this.content.coefficient)){
          return
        }        
        this.$store.dispatch('writeMaterial',this.content)
        this.content.dialog= false

      },
      close(){
        this.content.dialog = false
      }
    }
  }
</script>
