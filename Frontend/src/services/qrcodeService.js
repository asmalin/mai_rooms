import { saveAs } from "file-saver";
import QRCode from "https://esm.sh/qrcode";
import JSZip from "jszip";

export async function CreateQRCodes(roomIds) {
let roomsData = [];

let promises = roomIds.map(async function (roomId) {
    let roomName = await getRoomNameById(roomId);
    if (roomName != null) {
    roomsData.push({ id: roomId, name: roomName });
    }
});

await Promise.all(promises);

var opts = {
    errorCorrectionLevel: "H",
    type: "image/png",
    quality: 0.3,
    margin: 1,
    width: 300,
};

const zip = new JSZip();
promises = [];

for (const room of roomsData) {
    const promise = new Promise((resolve, reject) => {
    QRCode.toDataURL(
        website_domain + `/rooms?roomId=${room.id}`,
        opts,
        function (err, qrCodeImageDataUrl) {
        if (err) {
            reject(err);
        } else {
            downloadWithCaption(
            qrCodeImageDataUrl,
            `qrcode_${room.name}.png`,
            room.name
            )
            .then((blob) => {
                zip.file(`qrcode_${room.name}.png`, blob);
                resolve();
            })
            .catch((error) => {
                reject(error);
            });
        }
        }
    );
    });
    promises.push(promise);
}

Promise.all(promises)
    .then(() => {
    zip.generateAsync({ type: "blob" }).then(function (content) {
        saveAs(content, "qrcodes.zip");
    });
    })
    .catch((error) => {
    console.error("Error generating QR codes and creating zip:", error);
    });
}

function downloadWithCaption(dataurl, caption) {
return new Promise((resolve, reject) => {
    const downloadCanvas = document.createElement("canvas");
    downloadCanvas.width = 350;
    downloadCanvas.height = 350;
    const ctx = downloadCanvas.getContext("2d");

    const img = new Image();
    img.crossOrigin = "Anonymous";
    img.onload = function () {
    ctx.fillStyle = "white";
    ctx.fillRect(0, 0, downloadCanvas.width, downloadCanvas.height);
    ctx.fillStyle = "black";
    ctx.drawImage(img, 0, 0);

    ctx.font = "20px Times New Roman";
    ctx.fillText(caption, 70, 340);

    downloadCanvas.toBlob(function (blob) {
        if (blob) {
        resolve(blob);
        } else {
        reject(new Error("Failed to create blob."));
        }
    }, "image/jpeg");
    };
    img.onerror = function (error) {
    reject(error);
    };
    img.src = dataurl;
});
}

async function getRoomNameById(roomId) {
try {
    
    const response = await fetch("/api/room/" + roomId, {
    method: "GET",
    });

    if (!response.ok) {
    return null;
    }

    const userData = await response.json();
    return userData;
} catch (error) {
    return null;
}
}