//
//  AppDelegate.h
//  Squirrel
//
//  Created by JamesMao on 4/11/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import <UIKit/UIKit.h>

@class LoginViewController;

@interface AppDelegate : UIResponder <UIApplicationDelegate>
{
    LoginViewController* firstViewController_ ;
}

@property (strong, nonatomic) UIWindow *window;


@end

