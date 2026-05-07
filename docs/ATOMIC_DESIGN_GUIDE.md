# Atomic Design Guide

This project strictly follows Brad Frost's Atomic Design methodology.

## 1. Atoms
Atoms are the basic building blocks of matter. Applied to web interfaces, atoms are our HTML tags, such as a form label, an input or a button. They cannot be broken down any further without ceasing to be functional.
- **Location:** `src/components/atoms/`
- **Examples:** `BaseButton`, `BaseInput`, `BaseCard`, `BaseLoader`.

## 2. Molecules
Molecules are groups of atoms bonded together and are the smallest fundamental units of a compound. These take on their own properties and serve as the backbone of our design systems.
- **Location:** `src/components/molecules/`
- **Examples:** `LoginForm` (Input + Input + Button), `WalletCard` (Card + Typography + Button).

## 3. Organisms
Organisms are groups of molecules joined together to form a relatively complex, distinct section of an interface.
- **Location:** `src/components/organisms/`
- **Examples:** `SidebarNavigation`, `TransactionTable`, `DashboardHeader`.

## 4. Templates
Templates consist mostly of groups of organisms stitched together to form pages. It’s here where we start to see the design coming together and start seeing things like layout in action.
- **Location:** `src/components/templates/`
- **Examples:** `DashboardTemplate`, `AuthTemplate`.

## 5. Pages
Pages are specific instances of templates. Here, placeholder content is replaced with real representative content to give an accurate depiction of what a user will ultimately see. Pages are what the Vue Router points to.
- **Location:** `src/components/pages/`
- **Examples:** `LoginPage`, `DashboardPage`, `WalletPage`.
