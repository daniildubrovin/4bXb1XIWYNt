var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

// Get references to login form, signup form, and buttons to switch between them
const signinContainer = document.getElementById('signInContainer'); // Reference to login form
const signupContainer = document.getElementById('signUpContainer'); // Reference to signup form
const showSignupBtn = document.getElementById('showSignUp'); // Button to show signup form
const showSigninBtn = document.getElementById('showSignIn'); // Button to show login form

// Event listener to switch to signup form when the 'Sign Up' button is clicked
showSignupBtn.addEventListener('click', () => {
    signinContainer.classList.add('hidden'); // Hide login form
    signupContainer.classList.remove('hidden'); // Show signup form
	document.querySelectorAll('.error').forEach(e => e.remove());
});

// Event listener to switch to login form when the 'Sign In' button is clicked
showSigninBtn.addEventListener('click', () => {
    signinContainer.classList.remove('hidden'); // Show login form
    signupContainer.classList.add('hidden'); // Hide signup form
	document.querySelectorAll('.error').forEach(e => e.remove());
});