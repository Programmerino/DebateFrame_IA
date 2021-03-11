// Styles
require('codemirror/lib/codemirror.css');
require("github-markdown-css");
require("uikit/dist/css/uikit.min.css");
require("./styles.css");
require("tocbot/dist/tocbot.css");
require("medium-editor/dist/css/medium-editor.min.css");
require("medium-editor/dist/css/themes/beagle.min.css");
require("@fortawesome/fontawesome-free/css/all.css");

// HTML
require("./index.html");

// JavaScript
require("./go_wasm.js");
window.UIkit = require("uikit");
window.Dropzone = require("dropzone/dist/dropzone.js");
window.jQuery = require("jquery");
window.tocbot = require("tocbot");
window.MediumEditor = require("medium-editor");
window.rangy = require("rangy/lib/rangy-classapplier");
import Icons from 'uikit/dist/js/uikit-icons';
require("./vendor/fa-uikit/js/uikit-fa-icons"); // Basically styles as it just loads in the SVG files as icons



// loads the Icon plugin
UIkit.use(Icons);
if (!WebAssembly.instantiateStreaming) { // polyfill
    console.log("using pollyfill for instatiateStreaming")
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

const go = new Go();

WebAssembly.instantiateStreaming(fetch("client.wasm"), go.importObject).then((result) => {
    console.log("running go code")
    go.run(result.instance);
});

document.addEventListener("DOMContentLoaded", function (event) {
    UIkit.modal(document.getElementById("modal-loading")).show();
    rangy.init();

    window.HighlighterButton = MediumEditor.extensions.button.extend({
        name: 'highlighter',

        tagNames: ['mark'], // nodeName which indicates the button should be 'active' when isAlreadyApplied() is called
        contentDefault: '<b>H</b>', // default innerHTML of the button
        contentFA: '<i class="fa fa-paint-brush"></i>', // innerHTML of button when 'fontawesome' is being used
        aria: 'Highlight', // used as both aria-label and title attributes
        action: 'highlight', // used as the data-action attribute of the button

        init: function () {
            MediumEditor.extensions.button.prototype.init.call(this);

            this.classApplier = rangy.createClassApplier('highlight', {
                elementTagName: 'mark',
                normalize: true
            });
        },

        handleClick: function (event) {
            this.classApplier.toggleSelection();
            this.base.checkContentChanged();
        }
    });
});

function b64toBlob(b64Data, contentType) {
    var sliceSize = 512;

    var byteCharacters = atob(b64Data);
    var byteArrays = [];

    for (var offset = 0; offset < byteCharacters.length; offset += sliceSize) {
        var slice = byteCharacters.slice(offset, offset + sliceSize);

        var byteNumbers = new Array(slice.length);
        for (var i = 0; i < slice.length; i++) {
            byteNumbers[i] = slice.charCodeAt(i);
        }

        var byteArray = new Uint8Array(byteNumbers);

        byteArrays.push(byteArray);
    }

    var blob = new Blob(byteArrays, { type: contentType });
    return blob;
}

window.b64toBlob = b64toBlob


function blobToBytes(blob) {
    var fileReader = new FileReader();
    fileReader.onload = function (event) {
        console.log("Finished")
        var uint8View = new Uint8Array(event.target.result);
        window.result = uint8View
    };
    fileReader.readAsArrayBuffer(blob);
}

window.blobToBytes = blobToBytes

if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
        navigator.serviceWorker.register('/dist/service-worker.js').then(registration => {
            console.log('SW registered: ', registration);
        }).catch(registrationError => {
            console.log('SW registration failed: ', registrationError);
            alert("There was an issue setting up offline access: \"" + registrationError + "\". Please check your browsers privacy settings regarding Service Workers and try again. For more information, please check the Developer's Console. Keep in mind that this website will fail to work when your computer's internet turns off without this feature...")
        });
    });
}
