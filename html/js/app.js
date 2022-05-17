const vuetify = new Vuetify({
  theme: {
    themes: {
      dark: {
        primary: '#620ee',
        secondary: '#03dac5',
        accent: '#82B1FF',
        error: '#FF5252',
        info: '#2196F3',
        success: '#4CAF50',
        warning: '#FFC107',
        background: '#121212',
      },

    },
    options: {
      customProperties: true
    },
  },
});
const app = new Vue({
  el: "#app",
  vuetify: vuetify,
  data: () => ({
    search:"",
    host_result:null,
    maxPage:1,
    minPage:1,
    curPage:1,
    showLoading:false,

  }),
  created(){
    this.$vuetify.theme.dark = true;
  },
  mounted () {
    this.host_result={matches:[]}
  },
  methods: {
    shodanHostSearch(search,page,event){
      event.preventDefault();
      this.host_result={matches:[]}
      this.showLoading=true;
      var params = new URLSearchParams();
      params.append("search",search);
      console.log(params.get("search"));
      console.log("Search param submitted: "+params.get("search"))


      if(page!=-1&&page<0){
        page=1;
      }

      params.append("page",page);
      console.log("Page:",page);
      //If first time clicked on search, need page count
      if(page=="-1"){
        this.showLoading=true;
        axios({
          method: "post",
          url: "/shodanpagecount",
          data: params,
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
        }).then( (response)=>{
          var tmp = parseInt(response.data);
          console.log("Pages:",tmp);
          if(tmp<=0){
            alert("No hosts found")
          } 
          if(tmp!=NaN&&tmp>0){
            this.maxPage = tmp;
            this.minPage = 1;
            this.curPage = 1;
            params.set("page",this.curPage);
            this.showLoading=true;
            axios({
              method: "post",
              url: "/shodansearch",
              data: params,
              headers: { "Content-Type": "application/x-www-form-urlencoded" },
            }).then( (response)=>{
              this.host_result = response.data
            }).catch((e)=>{
              this.showLoading=false;
              console.log(e);
              if(e.response){
                console.log(e.response.data);
                var err = e.response.data.error;
                alert(err);
              } 
            }).then(()=>{
            this.showLoading=false;
            });
          }
        }).catch((e)=>{
          this.showLoading=false;
          console.log(e)   
          alert("Error, Try again later !!!");
        }).then(()=>{
        });

      }else{
        axios({
          method: "post",
          url: "/shodansearch",
          data: params,
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
        }).then( (response)=>{
          this.host_result = response.data
          console.log(response.data)
        }).catch((e)=>{
          this.showLoading=false;
          console.log(e)   
          if(e.response){
            console.log(e.response.data);
            var err = e.response.data.error;
            alert(err);
          }
        }).then(()=>{
          this.showLoading=false;
        });
      }
    },
    nextPage(){
      if(this.curPage<this.maxPage){
        this.curPage++;
        var params = new URLSearchParams();
        params.append("search",this.search);
        params.append("page",this.curPage)
        axios({
          method: "post",
          url: "/shodansearch",
          data: params,
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
        }).then( (response)=>{
          this.host_result = response.data
        }).catch((e)=>{
          this.curPage--;
          this.showLoading=false;
          console.log(e);
          if(e.response){
            console.log(e.response.data);
            var err = e.response.data.error;
            alert(err);
          }
        }).then(()=>{
          this.showLoading=false;
        });
      }else{
        alert("Max Pages Reached !!!");
      }
      
    },
    prevPage(){
      if(this.curPage>this.minPage){
        this.curPage--;
        var params = new URLSearchParams();
        params.append("search",this.search);
        params.append("page",this.curPage)
        axios({
          method: "post",
          url: "/shodansearch",
          data: params,
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
        }).then( (response)=>{
          this.host_result = response.data
        }).catch((e)=>{
          this.curPage++;
          this.showLoading=false;
          console.log(e);
          if(e.response){
            console.log(e.response.data);
            var err = e.response.data.error;
            alert(err);
          }
        }).then(()=>{
          this.showLoading=false;
        });
      }else{
        alert("Min Page Reached !!!");
      }
      
    }
  },
});
