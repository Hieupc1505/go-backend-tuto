package uploadimage

type UploadResult struct {
	Url   string `json:"url"`
	Thumb string `json:"thumb"`
}

type IUploadImage interface {
	Upload(base64 string) (UploadResult, error)
}

//Imgur json

// {
// 	"status": 200,
// 	"success": true,
// 	"data": {
// 	  "id": "JRBePDz",
// 	  "deletehash": "EvHVZkhJhdNClgY",
// 	  "account_id": null,
// 	  "account_url": null,
// 	  "ad_type": null,
// 	  "ad_url": null,
// 	  "title": "Simple upload",
// 	  "description": "This is a simple image upload in Imgur",
// 	  "name": "",
// 	  "type": "image/jpeg",
// 	  "width": 600,
// 	  "height": 750,
// 	  "size": 54757,
// 	  "views": 0,
// 	  "section": null,
// 	  "vote": null,
// 	  "bandwidth": 0,
// 	  "animated": false,
// 	  "favorite": false,
// 	  "in_gallery": false,
// 	  "in_most_viral": false,
// 	  "has_sound": false,
// 	  "is_ad": false,
// 	  "nsfw": null,
// 	  "link": "https://i.imgur.com/JRBePDz.jpeg",
// 	  "tags": [],
// 	  "datetime": 1708424380,
// 	  "mp4": "",
// 	  "hls": ""
// 	}
//   }

//Imgbb json
// {
// 	"data": {
// 	  "id": "2ndCYJK",
// 	  "title": "c1f64245afb2",
// 	  "url_viewer": "https://ibb.co/2ndCYJK",
// 	  "url": "https://i.ibb.co/w04Prt6/c1f64245afb2.gif",
// 	  "display_url": "https://i.ibb.co/98W13PY/c1f64245afb2.gif",
// 	  "width":"1",
// 	  "height":"1",
// 	  "size": "42",
// 	  "time": "1552042565",
// 	  "expiration":"0",
// 	  "image": {
// 		"filename": "c1f64245afb2.gif",
// 		"name": "c1f64245afb2",
// 		"mime": "image/gif",
// 		"extension": "gif",
// 		"url": "https://i.ibb.co/w04Prt6/c1f64245afb2.gif",
// 	  },
// 	  "thumb": {
// 		"filename": "c1f64245afb2.gif",
// 		"name": "c1f64245afb2",
// 		"mime": "image/gif",
// 		"extension": "gif",
// 		"url": "https://i.ibb.co/2ndCYJK/c1f64245afb2.gif",
// 	  },
// 	  "medium": {
// 		"filename": "c1f64245afb2.gif",
// 		"name": "c1f64245afb2",
// 		"mime": "image/gif",
// 		"extension": "gif",
// 		"url": "https://i.ibb.co/98W13PY/c1f64245afb2.gif",
// 	  },
// 	  "delete_url": "https://ibb.co/2ndCYJK/670a7e48ddcb85ac340c717a41047e5c"
// 	},
// 	"success": true,
// 	"status": 200
//   }
