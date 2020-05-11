const getPageTitle = (pathname, menuData) => {
  if (!menuData || menuData.length === 0) {
    return '';
  }
  let title;
  for (let i = 0; i < menuData.length; i++) {
    const menu = menuData[i];
    if (menu.children) {
      for (let j = 0; j < menu.children.length; j++) {
        const subMenu = menu.children[j];
        if (subMenu.Url === pathname) {
          title = subMenu.Name;
          break;
        }
      }
    }
  }
  return title;
};

export default getPageTitle;
