new Vue({
    el: "#app",
    data: {
        username: "",
        password: "",
        show: false,
        chooseLogin: "choosed",
        chooseSignup: "",
    },
    methods: {
        newUsername(value) {
            this.username = value;
            if (this.username === "" || this.password === "") {
                this.show = false;
                return;
            }
            this.show = true;
        },
        newPassword(value) {
            this.password = value;
            if (this.username === "" || this.password === "") {
                this.show = false;
                return;
            }
            this.show = true;
        },
    },
})