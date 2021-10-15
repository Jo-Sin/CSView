<template>
  <div class="container">
    <h2>Orders</h2>

    <!-- Search and filter elements -->
    <div class="form-group">
      <div class="row mt-5">
        <div class="col-auto">
          <label for="startDate" class="col-form-label">From:</label>
        </div>
        <div class="col">
          <input
            id="startDate"
            class="form-control"
            type="date"
            ref="startDate"
            @change="setStartDate"
          />
        </div>
        <div class="col-auto">
          <label for="endDate" class="col-form-label">To:</label>
        </div>
        <div class="col">
          <input
            id="endDate"
            class="form-control"
            type="date"
            ref="endDate"
            @change="setEndDate"
          />
        </div>
      </div>
      <div class="row mt-4">
        <div class="col-auto">
          <label for="searchText" class="col-form-label">Search Order:</label>
        </div>
        <div class="col">
          <input
            id="searchText"
            type="text"
            class="form-control"
            ref="searchText"
            @change="setSearch"
          />
        </div>
        <div class="col-auto">
          <button class="btn btn-outline-dark" @click="search">Search</button>
        </div>
      </div>
    </div>


    <!-- Pagination  -->
    <div class="row mt-5">
      <div class="col">
        <button class="btn btn-primary" @click="prevPage" :disabled="page <= 1">
          Prev
        </button>
      </div>
      <div class="col">Page: {{ page }} of {{ maxPages }}</div>
      <div class="col">
        <button
          class="btn btn-primary"
          @click="nextPage"
          :disabled="page >= maxPages"
        >
          Next
        </button>
      </div>
    </div>


    <!-- Data is shown here -->
    <table class="table">
      <thead>
        <tr>
          <th>Order name</th>
          <th>Customer Company</th>
          <th>Customer name</th>
          <th>Order date</th>
          <th>Delivered amount</th>
          <th>Total amount</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(order, index) in orders" :key="index">
          <td>{{ order["order-name"] }}</td>
          <td>{{ order["com-name"] }}</td>
          <td>{{ order["cust-name"] }}</td>
          <td>{{ order["order-date"] }}</td>
          <td>{{ order["delivered-amount"] }}</td>
          <td>{{ order["total-amount"] }}</td>
        </tr>
      </tbody>
    </table>


    <!-- Pagination again (Can be considered redundant since list is too small to scroll) -->
    <div class="row mt-5">
      <div class="col">
        <button class="btn btn-primary" @click="prevPage" :disabled="page <= 1">
          Prev
        </button>
      </div>
      <div class="col">Page: {{ page }} of {{ maxPages }}</div>
      <div class="col">
        <button
          class="btn btn-primary"
          @click="nextPage"
          :disabled="page >= maxPages"
        >
          Next
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "App",
  data() {
    return {
      orders: [],
      maxPages: 1,
      page: 1,
      query: "",
      startDate: new Date(2000, 0, 1, 0, 0, 0, 0),
      endDate: new Date(2000, 0, 1, 0, 0, 0, 0),
    };
  },
  mounted() {
    fetch("http://localhost:8000/count")
      .then((res) => res.json())
      .then((data) => {
        this.maxPages = data["page-count"];
        this.$refs.startDate.value = data["min-date"].substring(0, 10);
        var now = new Date(data["max-date"]);
        var day = ("0" + now.getDate()).slice(-2);
        var month = ("0" + (now.getMonth() + 1)).slice(-2);
        var today = now.getFullYear() + "-" + month + "-" + day;
        this.$refs.endDate.value = today;
        this.startDate = this.$refs.startDate.value;
        this.endDate = this.$refs.endDate.value;

        this.getData();
      })
      .catch((err) => console.log(err.message));
  },

  methods: {

    //Function that fetches the data
    getData() {
      let params =
        this.page.toString() +
        "/S" +
        this.query +
        "/" +
        this.startDate +
        "/" +
        this.endDate;
      console.log(params);

      fetch("http://localhost:8000/mongo-orders/" + params)
        .then((res) => res.json())
        .then((data) => {
          this.orders = data;
          return fetch("http://localhost:8000/post-orders/" + params);
        })
        .then((res) => res.json())
        .then((data) => {
          for (let i = 0; i < data.length; i++) {
            if (data[i]["delivered-amount"] > 0)
              this.orders[i]["delivered-amount"] =
                data[i]["delivered-amount"].toFixed(2);
            else this.orders[i]["delivered-amount"] = "-";
            this.orders[i]["total-amount"] = data[i]["total-amount"].toFixed(2);
            let tempDate = new Date(this.orders[i]["order-date"]);
            console.log(tempDate);
            this.orders[i]["order-date"] = tempDate.toLocaleString();
          }
          return fetch("http://localhost:8000/count");
        })
        .then((res) => res.json())
        .then((data) => (this.maxPages = data["page-count"]))
        .catch((err) => console.log(err.message));
    },


    prevPage() {
      this.page -= 1;
      this.getData();
    },
    nextPage() {
      this.page += 1;
      this.getData();
    },


    setStartDate() {
      this.startDate = this.$refs.startDate.value;
      if (this.startDate > this.endDate) {
        var now = new Date(this.startDate);
        var day = ("0" + now.getDate()).slice(-2);
        var month = ("0" + (now.getMonth() + 1)).slice(-2);
        var today = now.getFullYear() + "-" + month + "-" + day;
        this.$refs.endDate.value = today;
        this.endDate = this.$refs.endDate.value;
      }
    },
    setEndDate() {
      this.endDate = this.$refs.endDate.value;
      if (this.startDate > this.endDate) {
        var now = new Date(this.endDate);
        var day = ("0" + now.getDate()).slice(-2);
        var month = ("0" + (now.getMonth() + 1)).slice(-2);
        var today = now.getFullYear() + "-" + month + "-" + day;
        this.$refs.startDate.value = today;
        this.startDate = this.$refs.startDate.value;
      }
    },
    setSearch() {
      this.query = this.$refs.searchText.value;
    },

    search() {
      this.page = 1;
      this.getData();
    },
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
