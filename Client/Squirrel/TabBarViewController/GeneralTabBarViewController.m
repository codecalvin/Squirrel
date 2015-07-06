//
//  AppDelegate.h
//  Squirrel
//
//  Created by JamesMao on 4/11/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "GeneralTabBarViewController.h"
#import "Define.h"

@implementation GeneralTabBarViewController

- (id)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil
{
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if (self)
    {
        [self addViewController];
        [self loadViewControllers];
        
        self.viewControllers = viewControllers_;
        [self loadTabBarImage:self];
    }
    return self;
}

- (void)viewDidLoad
{
    [super viewDidLoad];
}

- (void)viewDidUnload
{
    [super viewDidUnload];
}

- (BOOL)shouldAutorotateToInterfaceOrientation:(UIInterfaceOrientation)interfaceOrientation
{
    return (interfaceOrientation == UIInterfaceOrientationPortrait);
}

- (void)addViewController
{
    if (viewControllers_ != nil)
    {
        return;
    }
    
    viewControllers_ = [[NSMutableArray alloc] init];
}

- (void)loadViewControllers
{

}

- (void)addSubViewController:(Class)subViewControllerClass
                  withTitle:(NSString*)title
{
    UIViewController* viewController = [[subViewControllerClass alloc] init];
    viewController.title = title;
    UINavigationController* navigationController = [[UINavigationController alloc] initWithRootViewController:viewController];
    [viewControllers_ addObject:navigationController];
}

- (void)addSubViewControllerInstance:(UIViewController*)subViewController
                           withTitle:(NSString*)title
{
    UIViewController* viewController = subViewController;
    viewController.title = title;
    UINavigationController* navigationController = [[UINavigationController alloc] initWithRootViewController:viewController];
    [viewControllers_ addObject:navigationController];
}

- (void)loadTabBarImage:(UITabBarController* )tabBarController
{
    int tabBarCount = (int)[self.tabBar.items count];
    int tabBarLastIndex = tabBarCount - 1;
    for (int index = 0; index < tabBarLastIndex; index++)
    {
        [self setTabBarItemImageAndTitle:index];
    }
    
    if ([self hasInformationTabBar])
    {
        [self setInformationTabBarItemImageAndTitle];
    }
    
    if ([self customizeBackgroundImage])
    {
        [self setBackgroundImage];
    }
    
    if ([self customizeTabItemTextColor])
    {
        [self setTabItemTextColor];
    }
}

- (void)setTabBarItemImageAndTitle:(int)index
{
    UITabBar *tabBar = self.tabBar;
    if (index >= [tabBar.items count])
    {
        return;
    }
    UITabBarItem *tabBarItem1 = [tabBar.items objectAtIndex:index];
    
    if ([self customizeTabBarName])
    {
        tabBarItem1.title = [self getTitle:index];
    }
    
    [self setTabBarItem:tabBarItem1
                  title:tabBarItem1.title
          selectedImage:[self getSelectedImage:index]
        unselectedImage:[self getUnselectedImage:index]];
}

- (void)setInformationTabBarItemImageAndTitle
{
    int lastIndex = (int)([self.tabBar.items count] - 1);
    UITabBarItem *tabBarItem1 = [self.tabBar.items objectAtIndex:lastIndex];
    
    int index = [self getTabBarMoreIndex];;
    if ([self customizeTabBarName])
    {
        tabBarItem1.title = [self getTitle:index];
    }
    
    [self setTabBarItem:tabBarItem1
                  title:tabBarItem1.title
          selectedImage:[self getSelectedImage:index]
        unselectedImage:[self getUnselectedImage:index]];
}

- (int)getTabBarMoreIndex
{
    return TAB_BAR_MORE_INDEX;
}

- (NSString*)getSelectedImage:(int)index
{
    //  to do
    return [self getUnselectedImage:index];
}

- (NSString*)getUnselectedImage:(int)index
{
    NSString* imageName = nil;
    switch (index)
    {
        case TabItemIndex1:
        {
            imageName = TAB_BAR_UNSELECTED_IMAGE_1;
        }
            break;
        case TabItemIndex2:
        {
            imageName = TAB_BAR_UNSELECTED_IMAGE_2;
        }
            break;
        case TabItemIndex3:
        {
            imageName = TAB_BAR_UNSELECTED_IMAGE_3;
        }
            break;
        case TabItemIndex4:
        {
            imageName = TAB_BAR_UNSELECTED_IMAGE_4;
        }
            break;
            
        default:
            break;
    }
    //imageName = nil;
    return imageName;
}

- (NSString*)getTitle:(int)index
{
    return nil;
}

- (void)setBackgroundImage
{
    UIImage* tabBarBackground = [UIImage imageNamed:TAB_BAR_BACKGROUND];
    [[UITabBar appearance] setBackgroundImage:tabBarBackground];
}

- (void)setTabItemTextColor
{
}

- (void)setTabBarItem:(UITabBarItem*)tabBarItem
                title:(NSString*)title
        selectedImage:(NSString*)selectedName
      unselectedImage:(NSString*)unselectedName
{

    //NSString* pureNameSelected = [selectedName stringByDeletingPathExtension];
    //NSString *fullPathSelected=[[NSBundle mainBundle] pathForResource:pureNameSelected ofType:@"png" inDirectory:@"Resources/Image"];
   // NSString *fullPathSelected=[[NSBundle mainBundle] pathForResource:pureNameSelected ofType:@"png"];
    
    UIImage * imageSelected = [UIImage imageNamed:selectedName];
    UIImage * imageUnSelected = [UIImage imageNamed:unselectedName];
    [tabBarItem initWithTitle:title
                        image:imageUnSelected
                selectedImage:imageSelected];
}

- (BOOL)customizeTabBarName
{
    return NO;
}

- (BOOL)hasInformationTabBar
{
    return NO;
}


- (BOOL)customizeBackgroundImage
{
    return YES;
}

- (BOOL)customizeTabItemTextColor
{
    return YES;
}


@end
