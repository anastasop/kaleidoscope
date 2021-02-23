import './kaleidoscope.css';
import { Application, Controller } from "stimulus";

const imageUrl = "/random";

class ImageController extends Controller {
    static targets = [ "img" ];
    static values = { filename: String };

    connect() {
        this.loadRandomImage();
    }

    loadRandomImage() {
        fetch(imageUrl).then(response => response.json()).then(json => {
            this.fileNameValue = json.filename;
            this.imgTarget.src = json.data;
            this.imgTarget.title = json.filename;
        });
    }

    select(e) {
        switch(e.button) {
        case 1:
            e.preventDefault();
            this.loadRandomImage();
            break;
        default:
            display(this.imgTarget.src);
        }
    }
}

function init() {
    const application = Application.start();
    application.register("image", ImageController);

    document.querySelector("#poster").addEventListener('click', showBoard);
}

function display(data) {
    document.getElementById("poster").style.visibility = 'visible';
    document.getElementById("main").style.visibility = 'hidden';
    document.getElementById("gallery").style.visibility = 'hidden';
    document.querySelector("#poster img").src = data;
}

function showBoard() {
    document.getElementById("poster").style.visibility = 'hidden';
   document.getElementById("main").style.visibility = 'visible';
    document.getElementById("gallery").style.visibility = 'visible';
    document.querySelector("#poster img").src = null;
}

document.addEventListener('DOMContentLoaded', (event) => {
    init();
});
