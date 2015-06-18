//
//  AppDelegate.h
//  Squirrel
//
//  Created by JamesMao on 4/11/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import <UIKit/UIKit.h>


//--------------------------------------------------------------------------------------------------



#define TAB_BAR_MORE_INDEX 3


//--------------------------------------------------------------------------------------------------

typedef enum
{
    TabItemIndex1,
    TabItemIndex2,
    TabItemIndex3,
    TabItemIndex4,
    TabItemIndex5,
}  TabItemIndex;

//--------------------------------------------------------------------------------------------------

@interface GeneralTabBarViewController : UITabBarController
{
    NSMutableArray* viewControllers_;
}

// utility
- (void)addSubViewController:(Class)subViewControllerClass
                  withTitle:(NSString*)title;
- (void)addSubViewControllerInstance:(UIViewController*)subViewController
                   withTitle:(NSString*)title;

// override it to customize
- (BOOL)customizeTabBarName;
- (BOOL)hasInformationTabBar;
- (BOOL)customizeBackgroundImage;
- (BOOL)customizeTabItemTextColor;

@end
