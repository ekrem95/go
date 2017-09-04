export const server = 'http://localhost:1323/';

export const post = (url, keys, values) =>
  new Promise((res, rej) => {
    if (keys.length == values.length) {
      let params = ``;

      for (var i = 0; i < keys.length; i++) {
        params += `${keys[i]}=${values[i]}&`;
      }

      params = params.substring(0, params.length - 1);

      const http = new XMLHttpRequest();

      // const params = `first=${first}&second=${second}`;
      http.open('POST', url, true);

      http.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');

      http.onreadystatechange = function () {
          if (http.readyState == 4 && http.status == 200) {
            res(JSON.parse(http.responseText));
          }
        };

      http.onerror = () => rej('Network Error');

      http.send(params);
    }
  });
