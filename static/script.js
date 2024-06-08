// Get the modal
var modal = document.getElementById('regWindow');
var submitButton = document.getElementById('reg-button');
var repeatPassword = false;
// When the user clicks anywhere outside of the modal, close it
window.onclick = function (event) {
  if (event.target == modal) {
    modal.style.display = "none";
  }
}

// document.getElementById("regForm").addEventListener('submit', function(event) {
//   var isValid = true;
//   var email = document.getElementById("emailReg").value.trim();
//   var psw = document.getElementById("pswReg").value;
//   var pswRepeat = document.getElementById("pswRepeatReg").value;
//   if (!ValidMail(email)) { //Если ошибка
//     document.getElementById("emailErrorReg").style.display = "block";
//     isValid = false;
//   } else {
//     document.getElementById("emailErrorReg").style.display = "none";
//   }

//   if (psw.length < 6 || psw.length > 64) {
//     document.getElementById("pswErrorReg").style.display = "block";
//     isValid = false;
//   } else {
//     document.getElementById("pswErrorReg").style.display = "none";
//   }

//   if (psw !== pswRepeat) {
//     document.getElementById("pswRepeatErrorReg").style.display = "block";
//     isValid = false;
//   } else {
//     document.getElementById("pswRepeatErrorReg").style.display = "none";
//   }

//   if (!isValid) {
//     event.preventDefault();
//   }

// })

function ValidMail() {
    var re = /^[\w-\.]+@[\w-]+\.[a-z]{2,4}$/i;
    var myMail = document.getElementById('email').value;
    var valid = re.test(myMail);

    return valid;
}
// function isPasswordMatch() {
//   var password = document.getElementById('psw').value;
//   var confirmPassword = document.getElementById('psw-repeat').value;
//   console.log(password.length)
//   repeatPassword = (password == confirmPassword);
//   if (password.length < 6) repeatPassword = false;
// }

// function registration() {
//   if (repeatPassword) submitButton.setAttribute("type", "submit");
//   else {
//     submitButton.setAttribute("type", "button");
//   }
// }


var modal1 = document.getElementById('authWindow');
var submitButton = document.getElementById('reg-button');
// When the user clicks anywhere outside of the modal, close it
window.onclick = function (event) {
  if (event.target == modal1) {
    modal1.style.display = "none";
  }
}


// document.getElementById("authForm").addEventListener('submit', function(event) {
//   var isValid = true;
//   var email = document.getElementById("email").value.trim();
//   var psw = document.getElementById("psw").value;
//   if (!ValidMail(email)) { //Если ошибка
//     document.getElementById("emailError").style.display = "block";
//     isValid = false;
//   } else {
//     document.getElementById("emailError").style.display = "none";
//   }

//   if (psw.length < 6 || psw.length > 64) {
//     document.getElementById("pswError").style.display = "block";
//     isValid = false;
//   } else {
//     document.getElementById("pswError").style.display = "none";
//   }

//   if (!isValid) {
//     event.preventDefault();
//   }

// })

document.addEventListener("DOMContentLoaded", function() {
  count_amount();
});


function open_menu() {
  let clicked = document.getElementById('drop-menu');
  if (clicked.style.display == 'block') {
      clicked.style.display = 'none';
  }
  else {
      clicked.style.display = 'block';
  }
}
// document.addEventListener("DOMContentLoaded", function() {
//   count_like_amount();
// });

// function count_like_amount() {
//   let count1 = document.getElementById("current_like_amount");
  

//   if (!count1) {
//     console.error("Element with id 'current_like_amount' not found.");
//     return;
//   }

//   let plusButton1 = document.getElementById("like");
//   let minusButton1 = document.getElementById("dislike");

//   if (!plusButton1 || !minusButton1) {
//     console.error("Buttons not found.");
//     return;
//   }

//   plusButton1.onclick = function() {
//     count1.innerText = parseInt(count1.innerText) + 1;
//     if (parseInt(count1.innerText) > 0){
//       count1.style.color = 'green'
//     } else if (parseInt(count1.innerText) < 0){
//       count1.style.color = 'red'
//     }else{
//       count1.style.color = 'black'
//     }
//   };

//   minusButton1.onclick = function() {
//     count1.innerText = parseInt(count1.innerText) - 1;
//     if (parseInt(count1.innerText) > 0){
//       count1.style.color = 'green'
//     } else if (parseInt(count1.innerText) < 0){
//       count1.style.color = 'red'
//     }else{
//       count1.style.color = 'black'
//     }
//   };
  
// }

document.addEventListener('DOMContentLoaded', function() {
  let commentFooters = document.querySelectorAll('.comment_footer');

  commentFooters.forEach(function(footer) {
    let countElement = footer.querySelector('.like_count');
    let likeButton = footer.querySelector('.like');
    let dislikeButton = footer.querySelector('.dislike');

    let liked = false;
    let disliked = false;

    likeButton.onclick = function() {
      likeButton.style.opacity = "1"
      dislikeButton.style.opacity = "0.7"
      likeButton.style.background = "#96e996"
      dislikeButton.style.background = "#eeeeee"
      if (!liked && !disliked) {
        //countElement.innerText = parseInt(countElement.innerText) + 1;
        liked = true;
        updateColor(countElement);
      } else if (disliked) {
        //countElement.innerText = parseInt(countElement.innerText) + 2;
        liked = true;
        disliked = false;
        updateColor(countElement);
      } else if(liked) {
        likeButton.style.opacity = "0.7"
        dislikeButton.style.opacity = "0.7"
        likeButton.style.background = "#eeeeee"
        dislikeButton.style.background = "#eeeeee"
        //countElement.innerText = parseInt(countElement.innerText) -1;
        liked = false;
        disliked = false;
        updateColor(countElement);
      }
    };

    dislikeButton.onclick = function() {
      likeButton.style.opacity = "0.7"
      dislikeButton.style.opacity = "1"
      likeButton.style.background = "#eeeeee"
      dislikeButton.style.background = "#e99696"
      if (!disliked && !liked) {
        //countElement.innerText = parseInt(countElement.innerText) - 1;
        disliked = true;
        updateColor(countElement);
      } else if (liked) {
        //countElement.innerText = parseInt(countElement.innerText) - 2;
        disliked = true;
        liked = false;
        updateColor(countElement);
      } else if(disliked) {
        likeButton.style.opacity = "0.7"
        dislikeButton.style.opacity = "0.7"
        likeButton.style.background = "#eeeeee"
        dislikeButton.style.background = "#eeeeee"
        //countElement.innerText = parseInt(countElement.innerText) +1;
        liked = false;
        disliked = false;
        updateColor(countElement);
      }
    };

    function updateColor(element) {
      if (parseInt(element.innerText) > 0) {
        element.style.color = 'green';
      } else if (parseInt(element.innerText) < 0) {
        element.style.color = 'red';
      } else {
        element.style.color = 'black';
      }
    }
  });
});



document.addEventListener("DOMContentLoaded", function() {
  const sizes = document.querySelectorAll('.size');
  
  sizes.forEach(size => {
    size.addEventListener('click', function() {
      // Удаление класса selected у всех кнопок
      sizes.forEach(s => s.classList.remove('selected'));
      // Добавление класса selected к нажатой кнопке
      this.classList.add('selected');
    });
  });
  
});

function count_amount() {
  let count = document.getElementById("current_amount");

  if (!count) {
    console.error("Element with id 'current_amount' not found.");
    return;
  }

  let plusButton = document.getElementById("plus");
  let minusButton = document.getElementById("minus");

  if (!plusButton || !minusButton) {
    console.error("Buttons not found.");
    return;
  }

  plusButton.onclick = function() {
    count.innerText = parseInt(count.innerText) + 1;
  };

  minusButton.onclick = function() {
    let currentValue = parseInt(count.innerText);
    if (currentValue > 0) {
      count.innerText = currentValue - 1;
    }
  };
}




