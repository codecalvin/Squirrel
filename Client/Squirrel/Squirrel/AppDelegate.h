//
//  AppDelegate.h
//  Squirrel
//
//  Created by JamesMao on 4/11/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import <UIKit/UIKit.h>

@interface AppDelegate : UIResponder <UIApplicationDelegate>
{
    UITabBarController *tabBarController_;
    UIViewController* firstViewController_;
}

@property (strong, nonatomic) UIWindow *window;


@end

